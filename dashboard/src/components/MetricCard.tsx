import React from 'react';
import { LucideIcon } from 'lucide-react';

interface MetricCardProps {
  title: string;
  value: string | number;
  unit?: string;
  icon: LucideIcon;
  trend?: {
    value: number;
    isUp: boolean;
  };
  children?: React.ReactNode;
}

const MetricCard: React.FC<MetricCardProps> = ({
  title,
  value,
  unit,
  icon: Icon,
  trend,
  children,
}) => {
  return (
    <div className="bg-card border border-border p-5 rounded-xl shadow-sm hover:border-accent transition-colors">
      <div className="flex justify-between items-start mb-4">
        <div className="p-2 bg-secondary rounded-lg">
          <Icon className="w-5 h-5 text-accent" />
        </div>
        {trend && (
          <div
            className={`text-xs font-medium ${
              trend.isUp ? 'text-red-400' : 'text-green-400'
            }`}
          >
            {trend.isUp ? '↑' : '↓'} {trend.value}%
          </div>
        )}
      </div>
      <div>
        <h3 className="text-muted-foreground text-sm font-medium mb-1">{title}</h3>
        <div className="flex items-baseline gap-1">
          <span className="text-2xl font-bold tracking-tight">{value}</span>
          {unit && <span className="text-muted-foreground text-sm">{unit}</span>}
        </div>
      </div>
      <div className="mt-4">{children}</div>
    </div>
  );
};

export default MetricCard;
