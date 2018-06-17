import threading
import paramiko
import logging
import time
import socket
cmd = "cat /proc/meminfo"
metrics={}
#add your hosts here
hosts = ["10.1.8.19","10.1.9.30",]
clients={}
host="127.0.0.1"
port=8125
username="" # dkorol
keyPath="" # /Users/dmitriy_korol/.ssh/id_rsa
class Metric(object):
    def __init__(self,  FreeMemory=0):
        self.lock = threading.Lock()
        self.FreeMemory=FreeMemory
    def updateMemory(self,amount):
        logging.debug('Waiting for lock')
        self.lock.acquire()
        try:
            logging.debug('Acquired lock')
            self.FreeMemory = amount
        finally:
            self.lock.release()
    def updateMemory(self,amount):
        logging.debug('Waiting for lock')
        self.lock.acquire()
        try:
            logging.debug('Acquired lock')
            self.FreeMemory = amount
        finally:
            self.lock.release()
    def readMemory(self):
        logging.debug('Waiting for lock')
        self.lock.acquire()
        try:
            logging.debug('Acquired lock')
            value = self.FreeMemory
            self.lock.release()
        finally:
            return value



def updater(host):
    while True:
        ssh = paramiko.SSHClient()
        ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())
        ssh.connect(host, username=username, password='xy',key_filename=keyPath)
        stdin, stdout, stderr = ssh.exec_command(cmd)
        stdin.write('xy\n')
        stdin.flush()
        metric=metrics[host]
        str = stdout.readlines()[1].split(':')[1]
        freeMemory = [int(s) for s in str.split() if s.isdigit()]
        metric.updateMemory(freeMemory[0]/1024.0)
        time.sleep(15)


def sendMetric(host):
    while True:
        #metric example ---> metricname:0.50|g
        value = "infra2_{0}:{1}|g\n".format(host, metrics[host].readMemory())
        clients[host].send(value.replace(".", "_"))
        print value
        time.sleep(30)



def main():
    threads = []
    threads2=[]
    for h in hosts:
        tcpClient = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        tcpClient.connect((host, port))
        clients[h]=tcpClient
        metrics[h]=Metric()
        t = threading.Thread(target=updater, args=(h,),)
        t2 = threading.Thread(target=sendMetric, args=(h,), )
        t.setDaemon(True)
        t2.setDaemon(True)
        t.start()
        t2.start()
        threads.append(t)
        threads2.append(t)
    while True:
        time.sleep(15)


main()
