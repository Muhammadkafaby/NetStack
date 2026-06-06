import { useState, useEffect, useRef } from 'react';
import { HostMetrics } from '../types/metrics';
import { generateMockMetrics } from '../mocks/metrics';

export function useMetrics(hostname?: string) {
  const [metrics, setMetrics] = useState<HostMetrics | HostMetrics[] | null>(null);
  const [isConnected, setIsConnected] = useState(false);
  const socketRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    // If we're in development or no backend URL is provided, use mock data
    const useMock = true; // In real app, this would be based on env or connection success

    if (useMock) {
      const interval = setInterval(() => {
        if (hostname) {
          setMetrics(generateMockMetrics(hostname));
        } else {
          // Update all hosts if no specific hostname
          setMetrics(prev => {
            if (Array.isArray(prev)) {
              return prev.map(m => ({
                ...m,
                cpu: { ...m.cpu, usage: Math.floor(Math.random() * 100) },
                memory: { ...m.memory, used: m.memory.used + (Math.random() - 0.5) * 100 * 1024 * 1024 }
              }));
            }
            return prev;
          });
        }
        setIsConnected(true);
      }, 2000);

      return () => clearInterval(interval);
    }

    // WebSocket implementation skeleton
    const connect = () => {
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      const wsUrl = `${protocol}//${window.location.host}/api/ws${hostname ? `/${hostname}` : ''}`;
      
      socketRef.current = new WebSocket(wsUrl);

      socketRef.current.onopen = () => {
        setIsConnected(true);
        console.log('Connected to VisiMon metrics stream');
      };

      socketRef.current.onmessage = (event) => {
        const data = JSON.parse(event.data);
        setMetrics(data);
      };

      socketRef.current.onclose = () => {
        setIsConnected(false);
        console.log('Disconnected from VisiMon metrics stream. Retrying...');
        setTimeout(connect, 5000);
      };

      socketRef.current.onerror = (error) => {
        console.error('WebSocket error:', error);
        socketRef.current?.close();
      };
    };

    connect();

    return () => {
      socketRef.current?.close();
    };
  }, [hostname]);

  return { metrics, isConnected };
}
