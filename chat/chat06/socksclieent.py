import socket
import struct
import time

def test_socks5_udp(socks5_host, socks5_port, target_host, target_port):
    # 1. 建立与 SOCKS5 代理的 TCP 连接（UDP 关联需要先握手）
    try:
        tcp_sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        tcp_sock.connect((socks5_host, socks5_port))
        print(f"✅ 已连接到 SOCKS5 代理 {socks5_host}:{socks5_port}")
    except Exception as e:
        print(f"❌ 连接 SOCKS5 代理失败：{e}")
        return

    # 2. SOCKS5 握手（认证阶段，无密码）
    handshake = struct.pack('!BBB', 0x05, 0x01, 0x00)
    tcp_sock.send(handshake)
    response = tcp_sock.recv(2)
    if len(response) != 2 or response[0] != 0x05 or response[1] != 0x00:
        print(f"❌ SOCKS5 握手失败，响应：{response.hex()}")
        tcp_sock.close()
        return
    print("✅ SOCKS5 握手成功（无需认证）")

    # 3. 请求 UDP 关联（核心步骤，测试 UDP 支持）
    udp_associate = struct.pack('!BBBBIH', 0x05, 0x03, 0x00, 0x01, 0, 0)
    tcp_sock.send(udp_associate)

    # 接收 UDP 关联响应
    response = tcp_sock.recv(10)
    if len(response) < 10 or response[0] != 0x05:
        print(f"❌ UDP 关联响应格式错误，响应：{response.hex()}")
        tcp_sock.close()
        return

    # 解析响应：状态码(第2字节)、代理分配的 UDP 端口
    status = response[1]
    if status != 0x00:
        print(f"❌ SOCKS5 代理不支持 UDP！状态码：{status}（0x00=成功，其他=失败）")
        tcp_sock.close()
        return

    # 提取代理分配的 UDP 端口（用于发送/接收 UDP 数据）
    proxy_udp_port = struct.unpack('!H', response[8:10])[0]
    print(f"✅ UDP 关联成功！代理分配的 UDP 端口：{proxy_udp_port}")

    # 4. 创建本地 UDP 套接字（用于和代理的 UDP 端口通信）
    local_udp_sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    # 绑定本地随机端口（避免端口占用）
    local_udp_sock.bind(('127.0.0.1', 0))
    print(f"🔌 本地 UDP 套接字已绑定：{local_udp_sock.getsockname()}")

    # 5. 构造 SOCKS5 UDP 数据包并发送
    # SOCKS5 UDP 头部格式：保留位(0x0000) + FRAG(0x00) + 地址类型(0x01) + 目标IP + 目标端口
    target_ip_int = struct.unpack('!I', socket.inet_aton(target_host))[0]
    udp_header = struct.pack('!BBBBIH', 0x00, 0x00, 0x00, 0x01, target_ip_int, target_port)
    send_data = b"test udp from socks5"
    udp_packet = udp_header + send_data

    # 发送数据到代理的 UDP 端口
    local_udp_sock.sendto(udp_packet, (socks5_host, proxy_udp_port))
    print(f"📤 已发送 UDP 数据到代理：{send_data.decode('utf-8')}")

    # 6. 监听并接收代理返回的 UDP 响应（设置超时，避免无限等待）
    local_udp_sock.settimeout(10)  # 10秒超时
    try:
        print("\n⌛ 等待接收 SOCKS5 代理返回的 UDP 响应...")
        response_packet, addr = local_udp_sock.recvfrom(1024)
        print(f"✅ 收到来自代理 {addr} 的 UDP 数据包（原始）：{response_packet.hex()}")

        # 解析 SOCKS5 UDP 响应头部，提取真实响应内容
        if len(response_packet) < 8:
            print("❌ 响应数据包格式错误：长度不足")
        else:
            # 跳过头部：保留位(2B) + FRAG(1B) + 地址类型(1B) + 目标IP(4B) + 目标端口(2B)
            header_len = 2 + 1 + 1 + 4 + 2  # 总计10字节
            real_response = response_packet[header_len:]
            print(f"🎉 解析后的 UDP 响应内容：{real_response.decode('utf-8', errors='ignore')}")

    except socket.timeout:
        print("❌ 超时未收到 UDP 响应（10秒）")
    except Exception as e:
        print(f"❌ 接收 UDP 响应失败：{e}")

    # 7. 关闭所有连接
    local_udp_sock.close()
    tcp_sock.close()
    print("\n🔚 所有连接已关闭")

if __name__ == "__main__":
    # 配置参数（根据你的环境调整）
    SOCKS5_HOST = "127.0.0.1"    # 代理地址
    SOCKS5_PORT = 1234           # 代理端口
    TARGET_HOST = "127.0.0.1"    # 本地 UDP 服务地址
    TARGET_PORT = 8848           # 本地 UDP 服务端口

    # 执行测试
    test_socks5_udp(SOCKS5_HOST, SOCKS5_PORT, TARGET_HOST, TARGET_PORT)