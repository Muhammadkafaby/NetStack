import React from 'react';
import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';

function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

interface StatusBadgeProps {
  status: 'online' | 'offline' | 'warning' | 'critical';
  className?: string;
}

const StatusBadge: React.FC<StatusBadgeProps> = ({ status, className }) => {
  const statusConfig = {
    online: { label: 'Online', color: 'bg-green-500/10 text-green-500 border-green-500/20' },
    offline: { label: 'Offline', color: 'bg-gray-500/10 text-gray-500 border-gray-500/20' },
    warning: { label: 'Warning', color: 'bg-yellow-500/10 text-yellow-500 border-yellow-500/20' },
    critical: { label: 'Critical', color: 'bg-red-500/10 text-red-500 border-red-500/20' },
  };

  const config = statusConfig[status];

  return (
    <span
      className={cn(
        'px-2 py-0.5 rounded-full text-xs font-medium border uppercase tracking-wider',
        config.color,
        className
      )}
    >
      {config.label}
    </span>
  );
};

export default StatusBadge;
