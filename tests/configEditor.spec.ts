import { test, expect } from '@grafana/plugin-e2e';

test('smoke: should render config editor with URL field', async ({
  createDataSourceConfigPage,
  readProvisionedDataSource,
  page,
}) => {
  const ds = await readProvisionedDataSource({ fileName: 'datasources.yml' });
  await createDataSourceConfigPage({ type: ds.type });
  await expect(page.getByRole('textbox', { name: 'URL' })).toBeVisible();
});

test('"Save & test" should be successful when Trino is reachable', async ({
  createDataSourceConfigPage,
  readProvisionedDataSource,
  page,
}) => {
  const ds = await readProvisionedDataSource({ fileName: 'datasources.yml' });
  const configPage = await createDataSourceConfigPage({ type: ds.type });
  await page.getByRole('textbox', { name: 'URL' }).fill(ds.url ?? 'http://trino:8080');
  await expect(configPage.saveAndTest()).toBeOK();
});

test('"Save & test" should fail when Trino URL is invalid', async ({
  createDataSourceConfigPage,
  readProvisionedDataSource,
  page,
}) => {
  const ds = await readProvisionedDataSource({ fileName: 'datasources.yml' });
  const configPage = await createDataSourceConfigPage({ type: ds.type });
  await page.getByRole('textbox', { name: 'URL' }).fill('http://invalid-host:9999');
  await expect(configPage.saveAndTest()).not.toBeOK();
});
