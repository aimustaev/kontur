import React, { useEffect, useState } from 'react';
import { Select } from 'antd';

interface WorkflowConfig {
  name: string;
  version: string;
  is_active: boolean;
  id: string;
}

interface ConfigSelectorProps {
  onConfigSelect?: (config: WorkflowConfig | null) => void;
}

export const ConfigSelector: React.FC<ConfigSelectorProps> = ({ onConfigSelect }) => {
  const [configs, setConfigs] = useState<WorkflowConfig[]>([]);
  const [selectedConfig, setSelectedConfig] = useState<string | null>(null);

  useEffect(() => {
    const fetchConfigs = async () => {
      try {
        const response = await fetch('/api-w/configs/summaries');
        if (!response.ok) {
          throw new Error('Failed to fetch configs');
        }
        const data = await response.json();
        setConfigs(data);
      } catch (error) {
        console.error('Error fetching configs:', error);
      }
    };

    fetchConfigs();
  }, []);

  const handleConfigChange = (value: string) => {
    setSelectedConfig(value);
    const [name, version] = value.split('-');
    const selectedConfigData = configs.find(
      config => config.name === name && config.version === version
    );
    onConfigSelect?.(selectedConfigData || null);
  };

  return (
    <Select
      style={{ width: 300 }}
      placeholder="Выберите конфигурацию"
      onChange={handleConfigChange}
      value={selectedConfig}
    >
      {configs.map((config) => (
        <Select.Option 
          key={`${config.name}-${config.version}`} 
          value={`${config.name}-${config.version}`}
        >
          {config.name} ({config.version})
          {config.is_active && (
            <span style={{ 
              marginLeft: '8px',
              color: '#52c41a',
              fontSize: '12px'
            }}>
              • Активная
            </span>
          )}
        </Select.Option>
      ))}
    </Select>
  );
}; 