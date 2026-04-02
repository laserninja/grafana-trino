import React from 'react';
import { CodeEditor, Combobox, ComboboxOption, InlineField, Stack } from '@grafana/ui';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from '../datasource';
import { TrinoDataSourceOptions, TrinoQuery, FormatOption } from '../types';

type Props = QueryEditorProps<DataSource, TrinoQuery, TrinoDataSourceOptions>;

const FORMAT_OPTIONS: Array<ComboboxOption<FormatOption>> = [
  { label: 'Table', value: 'table' },
  { label: 'Time Series', value: 'time_series' },
  { label: 'Logs', value: 'logs' },
];

export function QueryEditor({ query, onChange, onRunQuery }: Props) {
  const onRawSqlChange = (rawSql: string) => {
    onChange({ ...query, rawSql });
  };

  const onFormatChange = (option: ComboboxOption<FormatOption>) => {
    onChange({ ...query, format: option.value });
    onRunQuery();
  };

  return (
    <Stack direction="column" gap={1}>
      <InlineField label="Format" labelWidth={10}>
        <Combobox
          id="query-editor-format"
          options={FORMAT_OPTIONS}
          value={query.format}
          onChange={onFormatChange}
          width={20}
        />
      </InlineField>
      <CodeEditor
        language="sql"
        value={query.rawSql ?? ''}
        onBlur={onRawSqlChange}
        onSave={onRawSqlChange}
        showMiniMap={false}
        showLineNumbers={true}
        height="200px"
      />
    </Stack>
  );
}
