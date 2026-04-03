import React from 'react';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { ConnectionSettings, Auth, AdvancedHttpSettings, convertLegacyAuthProps, ConfigSection } from '@grafana/plugin-ui';
import { Field, Input, SecretInput, Switch } from '@grafana/ui';
import { TrinoDataSourceOptions, TrinoSecureJsonData } from '../types';

interface Props extends DataSourcePluginOptionsEditorProps<TrinoDataSourceOptions, TrinoSecureJsonData> {}

export function ConfigEditor(props: Props) {
  const { options, onOptionsChange } = props;
  const { jsonData, secureJsonData, secureJsonFields } = options;

  const onJsonDataChange = <K extends keyof TrinoDataSourceOptions>(key: K, value: TrinoDataSourceOptions[K]) => {
    onOptionsChange({ ...options, jsonData: { ...jsonData, [key]: value } });
  };

  const onSecureJsonDataChange = <K extends keyof TrinoSecureJsonData>(key: K, value: string) => {
    onOptionsChange({ ...options, secureJsonData: { ...secureJsonData, [key]: value } });
  };

  const onResetSecureJsonData = (key: keyof TrinoSecureJsonData) => {
    onOptionsChange({
      ...options,
      secureJsonFields: { ...secureJsonFields, [key]: false },
      secureJsonData: { ...secureJsonData, [key]: '' },
    });
  };

  return (
    <>
      <ConnectionSettings config={options} onChange={onOptionsChange} />
      <Auth
        {...convertLegacyAuthProps({
          config: options,
          onChange: onOptionsChange,
        })}
      />

      <ConfigSection title="Trino Settings">
        <Field label="Impersonate logged-in user" description="Set the Trino session user to the current Grafana user">
          <Switch
            id="trino-enable-impersonation"
            value={jsonData.enableImpersonation ?? false}
            onChange={(e) => onJsonDataChange('enableImpersonation', e.currentTarget.checked)}
          />
        </Field>
        <Field label="Access token" description="Static access token for Trino authentication">
          <SecretInput
            id="trino-access-token"
            value={secureJsonData?.accessToken ?? ''}
            isConfigured={secureJsonFields?.accessToken ?? false}
            onChange={(e) => onSecureJsonDataChange('accessToken', e.currentTarget.value)}
            onReset={() => onResetSecureJsonData('accessToken')}
            width={40}
          />
        </Field>
        <Field label="Roles" description="Authorization roles for catalogs (e.g., system:roleS;catalog1:roleA;catalog2:roleB)">
          <Input
            id="trino-roles"
            value={jsonData.roles ?? ''}
            onChange={(e) => onJsonDataChange('roles', e.currentTarget.value)}
            width={60}
          />
        </Field>
        <Field label="Client tags" description="Comma-separated list of tags to identify Trino resource groups">
          <Input
            id="trino-client-tags"
            value={jsonData.clientTags ?? ''}
            onChange={(e) => onJsonDataChange('clientTags', e.currentTarget.value)}
            width={60}
            placeholder="tag1,tag2,tag3"
          />
        </Field>
      </ConfigSection>

      <ConfigSection title="OAuth2 Trino Authentication" isCollapsible isInitiallyOpen={false}>
        <Field label="Token URL" description="OAuth2 token endpoint URL for client credentials flow">
          <Input
            id="trino-token-url"
            value={jsonData.tokenUrl ?? ''}
            onChange={(e) => onJsonDataChange('tokenUrl', e.currentTarget.value)}
            width={60}
          />
        </Field>
        <Field label="Client ID">
          <Input
            id="trino-client-id"
            value={jsonData.clientId ?? ''}
            onChange={(e) => onJsonDataChange('clientId', e.currentTarget.value)}
            width={60}
          />
        </Field>
        <Field label="Client secret">
          <SecretInput
            id="trino-client-secret"
            value={secureJsonData?.clientSecret ?? ''}
            isConfigured={secureJsonFields?.clientSecret ?? false}
            onChange={(e) => onSecureJsonDataChange('clientSecret', e.currentTarget.value)}
            onReset={() => onResetSecureJsonData('clientSecret')}
            width={60}
          />
        </Field>
        <Field label="Impersonation user" description="If set, overrides the Trino session user for OAuth requests">
          <Input
            id="trino-impersonation-user"
            value={jsonData.impersonationUser ?? ''}
            onChange={(e) => onJsonDataChange('impersonationUser', e.currentTarget.value)}
            width={60}
          />
        </Field>
      </ConfigSection>

      <AdvancedHttpSettings config={options} onChange={onOptionsChange} />
    </>
  );
}
