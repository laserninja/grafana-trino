import { DataSourceInstanceSettings, CoreApp, ScopedVars } from '@grafana/data';
import { DataSourceWithBackend, getTemplateSrv } from '@grafana/runtime';
import { map } from 'lodash';

import { TrinoQuery, TrinoDataSourceOptions, DEFAULT_QUERY } from './types';
import { TrinoVariableSupport } from './variable';

export class DataSource extends DataSourceWithBackend<TrinoQuery, TrinoDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<TrinoDataSourceOptions>) {
    super(instanceSettings);
    this.variables = new TrinoVariableSupport();
    this.annotations = {};
    this.interpolateQueryExpr = this.interpolateQueryExpr.bind(this);
  }

  getDefaultQuery(_: CoreApp): Partial<TrinoQuery> {
    return DEFAULT_QUERY;
  }

  applyTemplateVariables(query: TrinoQuery, scopedVars: ScopedVars) {
    return {
      ...query,
      rawSql: getTemplateSrv().replace(query.rawSql, scopedVars, this.interpolateQueryExpr),
    };
  }

  filterQuery(query: TrinoQuery): boolean {
    return !!query.rawSql;
  }

  /**
   * Custom interpolation for template variables in SQL queries.
   * - Single-value: escapes single quotes (value' → value'')
   * - Multi-value: comma-separated quoted list ('val1','val2')
   */
  interpolateQueryExpr(value: string | string[], variable: { multi: boolean; includeAll: boolean }): string | string[] {
    if (!variable.multi && !variable.includeAll) {
      return escapeLiteral(value as string);
    }

    if (typeof value === 'string') {
      return quoteLiteral(value);
    }

    return map(value, quoteLiteral).join(',');
  }
}

function quoteLiteral(value: string): string {
  return "'" + String(value).replace(/'/g, "''") + "'";
}

function escapeLiteral(value: string): string {
  return String(value).replace(/'/g, "''");
}
