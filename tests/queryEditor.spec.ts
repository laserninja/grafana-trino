import { test, expect } from '@grafana/plugin-e2e';

test('smoke: should render query editor with format selector', async ({
  panelEditPage,
  readProvisionedDataSource,
}) => {
  const ds = await readProvisionedDataSource({ fileName: 'datasources.yml' });
  await panelEditPage.datasource.set(ds.name);

  const queryRow = panelEditPage.getQueryEditorRow('A');
  // Format selector should be visible with default "Table"
  await expect(queryRow.getByLabel('Format')).toBeVisible();
});

test('SELECT 1 should return a result', async ({ panelEditPage, readProvisionedDataSource, page }) => {
  const ds = await readProvisionedDataSource({ fileName: 'datasources.yml' });
  await panelEditPage.datasource.set(ds.name);
  await panelEditPage.setVisualization('Table');

  const queryRow = panelEditPage.getQueryEditorRow('A');
  // Type into the Monaco code editor
  const editor = queryRow.getByRole('code').locator('.view-lines');
  await editor.click();
  await page.keyboard.type('SELECT 1 AS value');

  await expect(panelEditPage.refreshPanel()).toBeOK();
  await expect(panelEditPage.panel.data).toContainText(['1']);
});
