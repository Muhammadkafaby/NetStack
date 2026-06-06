export interface MetricHistory {
  timestamp: number;
  value: number;
}

export interface HostMetrics {
  hostname: string;
  ip: string;
  status: 'online' | 'offline' | 'warning';
  cpu: {
    usage: number;
    cores: number[];
    history: MetricHistory[];
  };
  memory: {
    total: number;
    used: number;
    free: number;
    cached: number;
    history: MetricHistory[];
  };
  disk: {
    total: number;
    used: number;
    free: number;
    io_read: number;
    io_write: number;
    history: MetricHistory[];
  };
  network: {
    inbound: number;
    outbound: number;
    history_in: MetricHistory[];
    history_out: MetricHistory[];
  };
  load: {
    avg1: number;
    avg5: number;
    avg15: number;
    history: MetricHistory[];
  };
  uptime: number;
}
