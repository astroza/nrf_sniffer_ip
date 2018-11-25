# (c) 2018 felipe@astroza.cl - See LICENSE

from socket import *
import time

def find_sniffer(timeout=6.0):
	sniffers = {}

	s = socket(AF_INET, SOCK_DGRAM)
        s.setsockopt(SOL_SOCKET, SO_BROADCAST, 1)
	s.bind(('0.0.0.0', 5311))
	start_time = time.time()
	elapsed_time = 0
	while elapsed_time < timeout:
		s.settimeout(timeout - elapsed_time)
		try:
			m = s.recvfrom(32)
			proto_port = m[0].split(":")
			if len(proto_port) > 1:
				sniffer = (proto_port[0], m[1][0], proto_port[1])
				sniffers[sniffer] = True
		except:
			pass
		elapsed_time = time.time() - start_time
	return sniffers.keys()


class NetSerialDevice:
	def __init__(self, device_tuple):
		proto, addr, port = device_tuple
		if proto == "TCP":
			self.s = socket(AF_INET, SOCK_STREAM)
			self.s.connect((addr, int(port)))
		else:
			raise Exception(proto + "is not a supported protocol")
		self.recv_buffer = ""
		self.recv_buffer_pos = 0

	def readByte(self, timeout=None):
		if self.recv_buffer_pos < len(self.recv_buffer):
			byte = self.recv_buffer[self.recv_buffer_pos]
			self.recv_buffer_pos = self.recv_buffer_pos + 1
			return byte
		if timeout != None:
			self.s.settimeout(timeout)
		try:
			self.recv_buffer = self.s.recv(4096)
			self.recv_buffer_pos = 0
			return self.readByte(timeout)
		except:
			return None

	def writeList(self, array):
		self.s.send(array)

if __name__ == "__main__":
	dev_addr = find_sniffer()[0]
	ser = NetSerialDevice(dev_addr)
	while True:
		print ser.readByte(),
