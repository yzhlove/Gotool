import socket

# 本地 UDP 回显服务：收到数据后，返回带标识的响应
def udp_echo_server(host='127.0.0.1', port=8848):
    sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    sock.bind((host, port))
    print(f"📡 本地 UDP 回显服务已启动，监听 {host}:{port}")

    try:
        while True:
            # 接收客户端数据
            data, addr = sock.recvfrom(1024)
            print(f"\n✅ 收到来自 {addr} 的数据：{data.decode('utf-8', errors='ignore')}")

            # 发送响应数据（模拟真实服务的返回）
            response = f"[ECHO] {data.decode('utf-8', errors='ignore')}".encode('utf-8')
            sock.sendto(response, addr)
            print(f"📤 已返回响应：{response.decode('utf-8')}")
    except KeyboardInterrupt:
        print("\n🛑 服务端已停止")
    finally:
        sock.close()

if __name__ == "__main__":
    udp_echo_server()