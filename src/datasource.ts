import { DataSourceInstanceSettings, CoreApp, ScopedVars } from '@grafana/data';
import { DataSourceWithBackend, getTemplateSrv } from '@grafana/runtime';

import { TrinoQuery, TrinoDataSourceOptions, DEFAULT_QUERY } from './types';

export class DataSource extends DataSourceWithBackend<TrinoQuery, TrinoDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<TrinoDataSourceOptions>) {
    super(instanceSettings);
  }

  getDefaultQuery(_: CoreApp): Partial<TrinoQuery> {
    return DEFAULT_QUERY;
  }

  applyTemplateVariables(query: TrinoQuery, scopedVars: ScopedVars) {
    return {
      ...query,
      rawSql: getTemplateSrv().replace(query.rawSql, scopedVars),
    };
  }

  filterQuery(query: TrinoQuery): boolean {
    return !!query.rawSql;
  }
}
