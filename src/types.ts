import { DataSourceJsonData } from '@grafana/data';
import { DataQuery } from '@grafana/schema';

export type FormatOption = 'table' | 'time_series' | 'logs';

export interface TrinoQuery extends DataQuery {
  rawSql: string;
  format: FormatOption;
}

export const DEFAULT_QUERY: Partial<TrinoQuery> = {
  rawSql: '',
  format: 'table',
};

/**
 * Options configured for each Trino DataSource instance.
 * Stored in jsonData (non-sensitive).
 */
export interface TrinoDataSourceOptions extends DataSourceJsonData {
  enableImpersonation?: boolean;
  tokenUrl?: string;
  clientId?: string;
  impersonationUser?: string;
  roles?: string;
  clientTags?: string;
}

/**
 * Secure values sent to the backend only, never exposed to the frontend.
 */
export interface TrinoSecureJsonData {
  accessToken?: string;
  clientSecret?: string;
}
