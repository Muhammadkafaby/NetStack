import React from 'react';
import { Doughnut } from 'react-chartjs-2';
import { Chart as ChartJS, ArcElement, Tooltip, Legend, type ChartOptions } from 'chart.js';

ChartJS.register(ArcElement, Tooltip, Legend);

interface GaugeChartProps {
  value: number;
  label: string;
  size?: number;
}

const GaugeChart: React.FC<GaugeChartProps> = ({ value, label, size = 120 }) => {
  const data = {
    datasets: [
      {
        data: [value, 100 - value],
        backgroundColor: [
          value > 80 ? '#ef4444' : value > 60 ? '#f59e0b' : '#0ea5e9',
          '#334155',
        ],
        borderWidth: 0,
        circumference: 270,
        rotation: 225,
        cutout: '80%',
      },
    ],
  };

  const options: ChartOptions<'doughnut'> = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      tooltip: { enabled: false },
      legend: { display: false },
    },
  };

  return (
    <div className="relative flex flex-col items-center justify-center" style={{ width: size, height: size }}>
      <Doughnut data={data} options={options} />
      <div className="absolute inset-0 flex flex-col items-center justify-center pt-2">
        <span className="text-xl font-bold leading-none">{value.toFixed(0)}%</span>
        <span className="text-[10px] text-muted-foreground font-medium uppercase tracking-wider">{label}</span>
      </div>
    </div>
  );
};

export default GaugeChart;
