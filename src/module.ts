import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './datasource';
import { ConfigEditor } from './components/ConfigEditor';
import { QueryEditor } from './components/QueryEditor';
import { TrinoQuery, TrinoDataSourceOptions } from './types';

export const plugin = new DataSourcePlugin<DataSource, TrinoQuery, TrinoDataSourceOptions>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
