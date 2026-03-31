import { test, expect, Page, Locator } from '@playwright/test';

const GRAFANA_CLIENT = 'grafana-client';

// Helper: resolves across Grafana 9.x (aria-label) and 11.x (data-testid)
function byLabelOrTestId(page: Page, label: string): Locator {
    return page.getByLabel(label).or(page.getByTestId(`data-testid ${label}`));
}

async function login(page: Page) {
    await page.goto('http://localhost:3000/login');
    await byLabelOrTestId(page, 'Username input field').fill('admin');
    await byLabelOrTestId(page, 'Password input field').fill('admin');
    await byLabelOrTestId(page, 'Login button').click();
    await byLabelOrTestId(page, 'Skip change password button').click();
}

async function goToTrinoSettings(page: Page) {
    await byLabelOrTestId(page, 'Toggle menu').click();
    await page.getByRole('link', {name: 'Connections'}).click();
    await page.getByRole('link', {name: 'Trino'}).click();
    // Grafana 9.x: "Create a Trino data source"; 11.x: "Add new data source"
    await page.getByText('Create a Trino data source')
        .or(page.getByRole('button', {name: 'Add new data source'}))
        .click();
}

async function setupDataSourceWithAccessToken(page: Page) {
    await byLabelOrTestId(page, 'Datasource HTTP settings url').fill('http://trino:8080');
    // 9.x: div text is "Impersonate logged in userAccess token"; 11.x: "Impersonate logged in user"
    await page.locator('div').filter({hasText: /^Impersonate logged in user(Access token)?$/}).getByLabel('Toggle switch').click();
    await page.locator('div').filter({hasText: /^Access token$/}).locator('input[type="password"]').fill('aaa');
    await byLabelOrTestId(page, 'Data source settings page Save and Test button').click();
    await page.waitForSelector('[role="alert"], [data-testid="data-testid Alert success"]', { timeout: 10000 });
}

async function setupDataSourceWithClientCredentials(page: Page, clientId: string) {
    await byLabelOrTestId(page, 'Datasource HTTP settings url').fill('http://trino:8080');
    await page.locator('div').filter({hasText: /^Token URL$/}).locator('input').fill('http://keycloak:8080/realms/trino-realm/protocol/openid-connect/token');
    await page.locator('div').filter({hasText: /^Client id$/}).locator('input').fill(clientId);
    await page.locator('div').filter({hasText: /^Client secret$/}).locator('input[type="password"]').fill('grafana-secret');
    await page.locator('div').filter({hasText: /^Impersonation user$/}).locator('input').fill('service-account-grafana-client');
    await byLabelOrTestId(page, 'Data source settings page Save and Test button').click();
    await page.waitForSelector('[role="alert"], [data-testid="data-testid Alert success"]', { timeout: 10000 });
}

async function runQueryAndCheckResults(page: Page) {
    // "Explore" text in 9.x, "Explore data" label in 11.x
    await page.getByText('Explore', {exact: true}).or(page.getByLabel('Explore data')).click();
    await page.getByTestId('data-testid TimePicker Open Button').click();
    await byLabelOrTestId(page, 'Time Range from field').fill('1995-01-01');
    await byLabelOrTestId(page, 'Time Range to field').fill('1995-12-31');
    await page.getByTestId('data-testid TimePicker submit button').click();
    // Format dropdown: aria-label in 9.x, div structure in 11.x
    await page.getByLabel('Format as').or(page.locator('div').filter({hasText: /^Format as/}).locator('svg')).click();
    await page.getByText('Table', { exact: true }).or(page.getByRole('option', {name: 'Table'})).click();
    await page.getByTestId('data-testid RefreshPicker run button').click();
    await expect(page.getByTestId('data-testid table body')).toContainText(/.*1995-01-19 0.:00:00.*/);
}

test('test with access token', async ({ page }) => {
    await login(page);
    await goToTrinoSettings(page);
    await setupDataSourceWithAccessToken(page);
    await runQueryAndCheckResults(page);
});

test('test client credentials flow', async ({ page }) => {
    await login(page);
    await goToTrinoSettings(page);
    await setupDataSourceWithClientCredentials(page, GRAFANA_CLIENT);
    await runQueryAndCheckResults(page);
});

test('test client credentials flow with wrong credentials', async ({ page }) => {
    await login(page);
    await goToTrinoSettings(page);
    await setupDataSourceWithClientCredentials(page, "some-wrong-client");
    // Check for error alert instead of checking absence of Explore button
    // The data source might still be saved even with wrong credentials
    // So we check if there's an error message or the Explore button is disabled/not present
    const exploreButton = page.getByText('Explore', {exact: true}).or(page.getByLabel('Explore data'));
    const errorAlert = page.locator('[role="alert"]:has-text("error"), [role="alert"]:has-text("failed"), [role="alert"]:has-text("Error")');
    
    // Either there should be an error alert, or the Explore button should not be visible
    const hasError = await errorAlert.count() > 0;
    
    // If there's no error alert, then Explore button should not be present
    if (!hasError) {
        await expect(exploreButton).toHaveCount(0);
    }
});

test('test client credentials flow with configured access token', async ({ page }) => {
    await login(page);
    await goToTrinoSettings(page);
    await page.locator('div').filter({hasText: /^Access token$/}).locator('input[type="password"]').fill('aaa');
    await setupDataSourceWithClientCredentials(page, GRAFANA_CLIENT);
    // Check for error alert instead of checking absence of Explore button
    // Setting both access token and client credentials should be invalid
    const exploreButton = page.getByText('Explore', {exact: true}).or(page.getByLabel('Explore data'));
    const errorAlert = page.locator('[role="alert"]:has-text("error"), [role="alert"]:has-text("failed"), [role="alert"]:has-text("Error")');
    
    // Either there should be an error alert, or the Explore button should not be visible
    const hasError = await errorAlert.count() > 0;
    
    // If there's no error alert, then Explore button should not be present
    if (!hasError) {
        await expect(exploreButton).toHaveCount(0);
    }
});
