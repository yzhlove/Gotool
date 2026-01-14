import socket
import struct
import time

def socks5_udp_associate(socks5_host, socks5_port):
    """建立 SOCKS5 UDP 关联，返回代理分配的 UDP 端口"""
    # 1. 建立 TCP 连接
    tcp_sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    try:
        tcp_sock.connect((socks5_host, socks5_port))
        print(f"✅ 已连接到 SOCKS5 代理 {socks5_host}:{socks5_port}")
    except Exception as e:
        print(f"❌ 连接 SOCKS5 代理失败：{e}")
        return None, None

    # 2. SOCKS5 握手（无认证）
    handshake = struct.pack('!BBB', 0x05, 0x01, 0x00)
    tcp_sock.send(handshake)
    response = tcp_sock.recv(2)
    if len(response) != 2 or response[0] != 0x05 or response[1] != 0x00:
        print(f"❌ SOCKS5 握手失败：响应 {response.hex()}")
        tcp_sock.close()
        return None, None
    print("✅ SOCKS5 握手成功（无需认证）")

    # 3. 请求 UDP 关联
    udp_associate = struct.pack('!BBBBIH', 0x05, 0x03, 0x00, 0x01, 0, 0)
    tcp_sock.send(udp_associate)
    response = tcp_sock.recv(10)
    if len(response) < 10 or response[0] != 0x05 or response[1] != 0x00:
        print(f"❌ UDP 关联失败：响应 {response.hex()}")
        tcp_sock.close()
        return None, None

    # 提取代理分配的 UDP 端口
    proxy_udp_port = struct.unpack('!H', response[8:10])[0]
    print(f"✅ UDP 关联成功，代理 UDP 端口：{proxy_udp_port}")
    return tcp_sock, proxy_udp_port

def send_dns_udp_query(socks5_host, proxy_udp_port, dns_server, dns_port, domain):
    """构造 DNS UDP 查询包，通过 SOCKS5 代理发送"""
    # 1. 生成 DNS 查询包（标准 DNS 格式）
    dns_id = 0x1234  # 随机 ID
    # 标志位：0x0100 = 递归查询
    flags = 0x0100
    qdcount = 1  # 1 个查询
    ancount = 0
    nscount = 0
    arcount = 0
    dns_header = struct.pack('!HHHHHH', dns_id, flags, qdcount, ancount, nscount, arcount)

    # 构造查询域名（例：example.com → 3example3com0）
    qname = b''
    for part in domain.split('.'):
        qname += struct.pack('B', len(part)) + part.encode('utf-8')
    qname += b'\x00'  # 域名结束符
    qtype = 1  # A 记录
    qclass = 1  # IN 类
    dns_query = qname + struct.pack('!HH', qtype, qclass)
    dns_packet = dns_header + dns_query

    # 2. 构造 SOCKS5 UDP 数据包头部
    target_ip_int = struct.unpack('!I', socket.inet_aton(dns_server))[0]
    socks5_udp_header = struct.pack('!BBBBIH', 0x00, 0x00, 0x00, 0x01, target_ip_int, dns_port)
    socks5_udp_packet = socks5_udp_header + dns_packet

    # 3. 发送 UDP 数据包到代理
    udp_sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    udp_sock.bind(('127.0.0.1', 0))  # 绑定本地随机端口
    udp_sock.sendto(socks5_udp_packet, (socks5_host, proxy_udp_port))
    print(f"✅ 已发送 DNS 查询：{domain} → {dns_server}:{dns_port}")

    # 4. 接收并解析响应
    udp_sock.settimeout(10)
    try:
        response, addr = udp_sock.recvfrom(1024)
        print(f"✅ 收到 DNS 响应（来自 {addr}）")

        # 跳过 SOCKS5 UDP 头部（前 10 字节）
        dns_response = response[10:]
        # 解析 DNS 响应头部
        resp_header = struct.unpack('!HHHHHH', dns_response[:12])
        resp_id, resp_flags, resp_qdcount, resp_ancount = resp_header[:4]

        if resp_ancount > 0:
            # 跳过查询部分，提取答案
            offset = 12 + len(dns_query)
            # 解析 A 记录
            ans_name = dns_response[offset:offset+2]
            ans_type = struct.unpack('!H', dns_response[offset+2:offset+4])[0]
            ans_class = struct.unpack('!H', dns_response[offset+4:offset+6])[0]
            ans_ttl = struct.unpack('!I', dns_response[offset+6:offset+10])[0]
            ans_len = struct.unpack('!H', dns_response[offset+10:offset+12])[0]
            ans_ip = socket.inet_ntoa(dns_response[offset+12:offset+12+ans_len])

            print(f"🎉 DNS 解析成功！{domain} → {ans_ip}")
        else:
            print("❌ DNS 解析失败：无答案记录")
    except socket.timeout:
        print("❌ 接收 DNS 响应超时（10秒）")
    except Exception as e:
        print(f"❌ 解析 DNS 响应失败：{e}")
    finally:
        udp_sock.close()

if __name__ == "__main__":
    # 配置参数（根据你的环境调整）
    SOCKS5_HOST = "127.0.0.1"
    SOCKS5_PORT = 1234
    DNS_SERVER = "8.8.8.8"  # Google DNS
    DNS_PORT = 53
    TEST_DOMAIN = "example.com"

    # 执行测试
    tcp_sock, proxy_udp_port = socks5_udp_associate(SOCKS5_HOST, SOCKS5_PORT)
    if proxy_udp_port:
        send_dns_udp_query(SOCKS5_HOST, proxy_udp_port, DNS_SERVER, DNS_PORT, TEST_DOMAIN)
        tcp_sock.close()
    print("\n🔚 测试完成")