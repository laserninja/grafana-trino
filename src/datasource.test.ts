import { DataSource } from './datasource';
import { TrinoDataSourceOptions } from './types';
import { DataSourceInstanceSettings } from '@grafana/data';

// Create a minimal instance of DataSource for testing
function createDataSource(): DataSource {
  const instanceSettings: DataSourceInstanceSettings<TrinoDataSourceOptions> = {
    id: 1,
    uid: 'test',
    type: 'trino-datasource',
    name: 'Trino',
    url: 'http://localhost:8080',
    access: 'proxy',
    readOnly: false,
    jsonData: {},
    meta: {} as any,
  };
  return new DataSource(instanceSettings);
}

describe('DataSource', () => {
  let ds: DataSource;

  beforeEach(() => {
    ds = createDataSource();
  });

  describe('filterQuery', () => {
    it('should return false for empty rawSql', () => {
      expect(ds.filterQuery({ refId: 'A', rawSql: '', format: 'table' })).toBe(false);
    });

    it('should return true for non-empty rawSql', () => {
      expect(ds.filterQuery({ refId: 'A', rawSql: 'SELECT 1', format: 'table' })).toBe(true);
    });
  });

  describe('interpolateQueryExpr', () => {
    it('should escape single quotes for single-value variables', () => {
      const result = ds.interpolateQueryExpr("it's a test", { multi: false, includeAll: false });
      expect(result).toBe("it''s a test");
    });

    it('should quote single string for multi-value variables', () => {
      const result = ds.interpolateQueryExpr('value1', { multi: true, includeAll: false });
      expect(result).toBe("'value1'");
    });

    it('should join array values with comma for multi-value variables', () => {
      const result = ds.interpolateQueryExpr(['val1', 'val2', 'val3'], { multi: true, includeAll: false });
      expect(result).toBe("'val1','val2','val3'");
    });

    it('should escape quotes in multi-value array items', () => {
      const result = ds.interpolateQueryExpr(["it's", "a'test"], { multi: true, includeAll: false });
      expect(result).toBe("'it''s','a''test'");
    });

    it('should handle includeAll the same as multi', () => {
      const result = ds.interpolateQueryExpr(['a', 'b'], { multi: false, includeAll: true });
      expect(result).toBe("'a','b'");
    });
  });
});
