import React from 'react';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { ConnectionSettings, Auth, convertLegacyAuthProps } from '@grafana/plugin-ui';
import { TrinoDataSourceOptions } from '../types';

interface Props extends DataSourcePluginOptionsEditorProps<TrinoDataSourceOptions> {}

export function ConfigEditor(props: Props) {
  const { options, onOptionsChange } = props;

  return (
    <>
      <ConnectionSettings config={options} onChange={onOptionsChange} />
      <Auth
        {...convertLegacyAuthProps({
          config: options,
          onChange: onOptionsChange,
        })}
      />
    </>
  );
}
