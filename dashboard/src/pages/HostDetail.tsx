import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { ArrowLeft, Cpu, Database, HardDrive, Network, Clock, Info } from 'lucide-react';
import { generateMockMetrics } from '../mocks/metrics';
import { HostMetrics } from '../types/metrics';
import TimeSeriesChart from '../components/TimeSeriesChart';
import GaugeChart from '../components/GaugeChart';
import StatusBadge from '../components/StatusBadge';

const HostDetail: React.FC = () => {
  const { hostname } = useParams<{ hostname: string }>();
  const navigate = useNavigate();
  const [metrics, setMetrics] = useState<HostMetrics | null>(null);

  useEffect(() => {
    if (hostname) {
      setMetrics(generateMockMetrics(hostname));
      const interval = setInterval(() => {
        setMetrics(prev => {
          if (!prev) return null;
          return {
            ...prev,
            cpu: { 
              ...prev.cpu, 
              usage: Math.floor(Math.random() * 100),
              history: [...prev.cpu.history.slice(1), { timestamp: Date.now(), value: Math.floor(Math.random() * 100) }]
            },
            memory: { 
              ...prev.memory, 
              used: prev.memory.used + (Math.random() - 0.5) * 50 * 1024 * 1024,
              history: [...prev.memory.history.slice(1), { timestamp: Date.now(), value: Math.floor(Math.random() * 100) }]
            }
          };
        });
      }, 1000);
      return () => clearInterval(interval);
    }
  }, [hostname]);

  if (!metrics) return <div>Loading...</div>;

  const cpuChartData = {
    labels: metrics.cpu.history.map(h => new Date(h.timestamp).toLocaleTimeString()),
    datasets: [{
      label: 'CPU Usage (%)',
      data: metrics.cpu.history.map(h => h.value),
      borderColor: '#0ea5e9',
      backgroundColor: 'rgba(14, 165, 233, 0.1)',
      fill: true,
    }]
  };

  const memChartData = {
    labels: metrics.memory.history.map(h => new Date(h.timestamp).toLocaleTimeString()),
    datasets: [{
      label: 'Memory Usage (%)',
      data: metrics.memory.history.map(h => h.value),
      borderColor: '#10b981',
      backgroundColor: 'rgba(16, 185, 129, 0.1)',
      fill: true,
    }]
  };

  return (
    <div className="space-y-8">
      <div className="flex items-center gap-4">
        <button 
          onClick={() => navigate('/')}
          className="p-2 hover:bg-secondary rounded-lg transition-colors"
        >
          <ArrowLeft className="w-5 h-5" />
        </button>
        <div>
          <div className="flex items-center gap-3">
            <h1 className="text-3xl font-bold tracking-tight">{metrics.hostname}</h1>
            <StatusBadge status={metrics.status} />
          </div>
          <p className="text-muted-foreground font-mono text-sm">{metrics.ip}</p>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
        <div className="lg:col-span-1 space-y-6">
          <div className="bg-card border border-border rounded-xl p-6">
            <h3 className="font-bold flex items-center gap-2 mb-4">
              <Info className="w-4 h-4 text-accent" /> System Info
            </h3>
            <div className="space-y-4 text-sm">
              <div className="flex justify-between">
                <span className="text-muted-foreground">OS</span>
                <span className="font-medium">Ubuntu 22.04 LTS</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Kernel</span>
                <span className="font-medium">5.15.0-x86_64</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Uptime</span>
                <span className="font-medium">14 days, 2:43:12</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">CPU Cores</span>
                <span className="font-medium">{metrics.cpu.cores.length} Cores</span>
              </div>
            </div>
          </div>

          <div className="bg-card border border-border rounded-xl p-6 flex justify-around">
            <GaugeChart value={metrics.cpu.usage} label="CPU" />
            <GaugeChart value={(metrics.memory.used / metrics.memory.total) * 100} label="RAM" />
          </div>
        </div>

        <div className="lg:col-span-3 space-y-6">
          <div className="bg-card border border-border rounded-xl p-6">
            <div className="flex items-center gap-2 mb-6">
              <Cpu className="w-5 h-5 text-primary" />
              <h3 className="font-bold">CPU Performance</h3>
            </div>
            <TimeSeriesChart data={cpuChartData} height={250} showGrid={true} />
          </div>

          <div className="bg-card border border-border rounded-xl p-6">
            <div className="flex items-center gap-2 mb-6">
              <Database className="w-5 h-5 text-green-500" />
              <h3 className="font-bold">Memory Utilization</h3>
            </div>
            <TimeSeriesChart data={memChartData} height={250} showGrid={true} />
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div className="bg-card border border-border rounded-xl p-6">
          <div className="flex items-center gap-2 mb-6">
            <HardDrive className="w-5 h-5 text-yellow-500" />
            <h3 className="font-bold">Disk Usage & I/O</h3>
          </div>
          <div className="grid grid-cols-2 gap-4 mb-6">
            <div className="p-4 bg-secondary/30 rounded-lg">
              <p className="text-xs text-muted-foreground uppercase font-bold mb-1">Total Disk</p>
              <p className="text-xl font-bold">{(metrics.disk.total / (1024**3)).toFixed(0)} GB</p>
            </div>
            <div className="p-4 bg-secondary/30 rounded-lg">
              <p className="text-xs text-muted-foreground uppercase font-bold mb-1">Used</p>
              <p className="text-xl font-bold">{(metrics.disk.used / (1024**3)).toFixed(0)} GB</p>
            </div>
          </div>
          <div className="w-full bg-secondary h-2 rounded-full overflow-hidden">
            <div 
              className="h-full bg-yellow-500"
              style={{ width: `${(metrics.disk.used / metrics.disk.total) * 100}%` }}
            ></div>
          </div>
        </div>

        <div className="bg-card border border-border rounded-xl p-6">
          <div className="flex items-center gap-2 mb-6">
            <Network className="w-5 h-5 text-purple-500" />
            <h3 className="font-bold">Network Traffic</h3>
          </div>
          <div className="flex justify-between items-center h-full pb-8">
            <div className="text-center flex-1">
              <p className="text-3xl font-bold text-primary">{(metrics.network.inbound / 1024).toFixed(1)}</p>
              <p className="text-xs text-muted-foreground uppercase">Incoming (Mbps)</p>
            </div>
            <div className="w-px h-12 bg-border"></div>
            <div className="text-center flex-1">
              <p className="text-3xl font-bold text-accent">{(metrics.network.outbound / 1024).toFixed(1)}</p>
              <p className="text-xs text-muted-foreground uppercase">Outgoing (Mbps)</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default HostDetail;
