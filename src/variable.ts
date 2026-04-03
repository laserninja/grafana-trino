import { StandardVariableQuery, StandardVariableSupport } from '@grafana/data';
import { DataSource } from './datasource';
import { TrinoQuery } from './types';

export class TrinoVariableSupport extends StandardVariableSupport<DataSource> {
  toDataQuery(query: StandardVariableQuery): TrinoQuery {
    return {
      refId: 'TrinoVariableQuery',
      rawSql: query.query,
      format: 'table',
    };
  }
}
