import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Cpu, Database, HardDrive, Network } from 'lucide-react';
import MetricCard from '../components/MetricCard';
import TimeSeriesChart from '../components/TimeSeriesChart';
import { generateMockMetrics, mockHosts } from '../mocks/metrics';
import { HostMetrics } from '../types/metrics';

const Overview: React.FC = () => {
  const navigate = useNavigate();
  const [metrics, setMetrics] = useState<HostMetrics[]>([]);

  useEffect(() => {
    // Simulate fetching data for all hosts
    const data = mockHosts.map(host => generateMockMetrics(host));
    setMetrics(data);

    const interval = setInterval(() => {
      setMetrics(prev => prev.map(m => ({
        ...m,
        cpu: { ...m.cpu, usage: Math.floor(Math.random() * 100) },
        memory: { ...m.memory, used: m.memory.used + (Math.random() - 0.5) * 100 * 1024 * 1024 }
      })));
    }, 2000);

    return () => clearInterval(interval);
  }, []);

  const totalCpuUsage = metrics.reduce((acc, m) => acc + m.cpu.usage, 0) / (metrics.length || 1);
  const totalMemUsed = metrics.reduce((acc, m) => acc + m.memory.used, 0);
  const totalMemTotal = metrics.reduce((acc, m) => acc + m.memory.total, 0);
  const memPercent = (totalMemUsed / totalMemTotal) * 100;

  const chartData = {
    labels: Array.from({ length: 20 }).map((_, i) => `${i}s`),
    datasets: [
      {
        label: 'System Load',
        data: Array.from({ length: 20 }).map(() => Math.random() * 2),
        borderColor: '#0ea5e9',
        backgroundColor: 'rgba(14, 165, 233, 0.1)',
        fill: true,
      },
    ],
  };

  return (
    <div className="space-y-8">
      <div className="flex justify-between items-end">
        <div>
          <h1 className="text-3xl font-bold tracking-tight mb-2">Fleet Overview</h1>
          <p className="text-muted-foreground">Real-time performance metrics across {mockHosts.length} active hosts.</p>
        </div>
        <div className="flex gap-3">
          <div className="bg-card border border-border px-4 py-2 rounded-lg text-sm font-medium">
            Avg CPU: <span className="text-accent">{totalCpuUsage.toFixed(1)}%</span>
          </div>
          <div className="bg-card border border-border px-4 py-2 rounded-lg text-sm font-medium">
            Avg RAM: <span className="text-accent">{memPercent.toFixed(1)}%</span>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <MetricCard title="CPU Usage" value={totalCpuUsage.toFixed(1)} unit="%" icon={Cpu} trend={{ value: 2.5, isUp: true }}>
          <div className="w-full bg-secondary h-2 rounded-full overflow-hidden">
            <div 
              className={`h-full transition-all duration-500 ${totalCpuUsage > 80 ? 'bg-red-500' : totalCpuUsage > 60 ? 'bg-yellow-500' : 'bg-primary'}`}
              style={{ width: `${totalCpuUsage}%` }}
            ></div>
          </div>
        </MetricCard>
        
        <MetricCard title="Memory Usage" value={(totalMemUsed / (1024**3)).toFixed(1)} unit="GB" icon={Database}>
           <div className="w-full bg-secondary h-2 rounded-full overflow-hidden">
            <div 
              className="h-full bg-primary transition-all duration-500"
              style={{ width: `${memPercent}%` }}
            ></div>
          </div>
        </MetricCard>

        <MetricCard title="Disk I/O" value="12.4" unit="MB/s" icon={HardDrive} trend={{ value: 12, isUp: false }} />
        <MetricCard title="Network" value="842" unit="Kbps" icon={Network} />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2 bg-card border border-border rounded-xl p-6">
          <div className="flex justify-between items-center mb-6">
            <h3 className="font-bold">Aggregated System Load</h3>
            <div className="flex gap-2">
              <span className="flex items-center gap-1 text-xs text-muted-foreground">
                <span className="w-2 h-2 rounded-full bg-primary"></span> Load Avg
              </span>
            </div>
          </div>
          <TimeSeriesChart data={chartData} height={300} showGrid={true} />
        </div>

        <div className="bg-card border border-border rounded-xl p-6">
          <h3 className="font-bold mb-6">Active Hosts</h3>
          <div className="space-y-4">
            {metrics.map((host) => (
              <div 
                key={host.hostname} 
                onClick={() => navigate(`/host/${host.hostname}`)}
                className="flex items-center justify-between p-3 bg-secondary/30 rounded-lg border border-transparent hover:border-accent hover:bg-secondary/50 transition-all cursor-pointer"
              >
                <div className="flex items-center gap-3">
                  <div className={`w-2 h-2 rounded-full ${host.status === 'online' ? 'bg-green-500' : 'bg-yellow-500'}`}></div>
                  <div>
                    <p className="text-sm font-medium">{host.hostname}</p>
                    <p className="text-[10px] text-muted-foreground uppercase">{host.ip}</p>
                  </div>
                </div>
                <div className="text-right">
                  <p className="text-xs font-bold">{host.cpu.usage}%</p>
                  <p className="text-[10px] text-muted-foreground">CPU</p>
                </div>
              </div>
            ))}
          </div>
          <button className="w-full mt-6 py-2 text-sm font-medium text-accent hover:bg-accent/10 rounded-lg transition-colors">
            View All Hosts
          </button>
        </div>
      </div>
    </div>
  );
};

export default Overview;
