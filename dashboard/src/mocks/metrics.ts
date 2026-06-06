import { HostMetrics } from '../types/metrics';

const generateHistory = (count: number, min: number, max: number) => {
  const now = Date.now();
  return Array.from({ length: count }).map((_, i) => ({
    timestamp: now - (count - i) * 1000,
    value: Math.floor(Math.random() * (max - min + 1)) + min,
  }));
};

export const generateMockMetrics = (hostname: string): HostMetrics => {
  return {
    hostname,
    ip: '192.168.1.' + Math.floor(Math.random() * 255),
    status: Math.random() > 0.1 ? 'online' : 'warning',
    cpu: {
      usage: Math.floor(Math.random() * 100),
      cores: [Math.random() * 100, Math.random() * 100, Math.random() * 100, Math.random() * 100],
      history: generateHistory(60, 10, 90),
    },
    memory: {
      total: 16 * 1024 * 1024 * 1024,
      used: 8 * 1024 * 1024 * 1024 + Math.random() * 4 * 1024 * 1024 * 1024,
      free: 2 * 1024 * 1024 * 1024,
      cached: 6 * 1024 * 1024 * 1024,
      history: generateHistory(60, 40, 80),
    },
    disk: {
      total: 512 * 1024 * 1024 * 1024,
      used: 200 * 1024 * 1024 * 1024,
      free: 312 * 1024 * 1024 * 1024,
      io_read: Math.random() * 100,
      io_write: Math.random() * 50,
      history: generateHistory(60, 5, 20),
    },
    network: {
      inbound: Math.random() * 1000,
      outbound: Math.random() * 500,
      history_in: generateHistory(60, 100, 1000),
      history_out: generateHistory(60, 50, 500),
    },
    load: {
      avg1: Math.random() * 2,
      avg5: Math.random() * 1.5,
      avg15: Math.random() * 1.2,
      history: generateHistory(60, 0, 2),
    },
    uptime: 123456,
  };
};

export const mockHosts = [
  'server-production-01',
  'server-production-02',
  'db-master',
  'db-replica-01',
  'redis-cache',
];
