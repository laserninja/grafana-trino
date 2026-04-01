import { test, expect, Page, Locator } from '@playwright/test';

const GRAFANA_CLIENT = 'grafana-client';
const GRAFANA_VERSION = process.env.GRAFANA_VERSION || '9.5.0';
const GRAFANA_MAJOR = parseInt(GRAFANA_VERSION.split('.')[0], 10);
const isNewGrafana = GRAFANA_MAJOR >= 10;

// Grafana 9.x uses aria-label, 10+ uses data-testid
function sel(page: Page, name: string): Locator {
    return isNewGrafana
        ? page.getByTestId(`data-testid ${name}`)
        : page.getByLabel(name);
}

async function login(page: Page) {
    await page.goto('http://localhost:3000/login');
    await sel(page, 'Username input field').fill('admin');
    await sel(page, 'Password input field').fill('admin');
    await sel(page, 'Login button').click();
    await sel(page, 'Skip change password button').click();
    await page.waitForLoadState('networkidle');
}

async function goToTrinoSettings(page: Page) {
    await sel(page, 'Toggle menu').click();
    await page.getByRole('link', {name: 'Connections'}).click();
    await page.getByRole('link', {name: 'Trino'}).click();
    if (isNewGrafana) {
        await page.getByRole('button', {name: 'Add new data source'}).click();
    } else {
        await page.getByText('Create a Trino data source').click();
    }
}

async function setupDataSourceWithAccessToken(page: Page) {
    await sel(page, 'Datasource HTTP settings url').fill('http://trino:8080');
    if (isNewGrafana) {
        await page.locator('div').filter({hasText: /^Impersonate logged in user$/}).getByLabel('Toggle switch').click();
    } else {
        await page.locator('div').filter({hasText: /^Impersonate logged in userAccess token$/}).getByLabel('Toggle switch').click();
    }
    await page.locator('div').filter({hasText: /^Access token$/}).locator('input[type="password"]').fill('aaa');
    await sel(page, 'Data source settings page Save and Test button').click();
    await page.waitForSelector('[role="alert"]', { timeout: 10000 });
}

async function setupDataSourceWithClientCredentials(page: Page, clientId: string) {
    await sel(page, 'Datasource HTTP settings url').fill('http://trino:8080');
    await page.locator('div').filter({hasText: /^Token URL$/}).locator('input').fill('http://keycloak:8080/realms/trino-realm/protocol/openid-connect/token');
    await page.locator('div').filter({hasText: /^Client id$/}).locator('input').fill(clientId);
    await page.locator('div').filter({hasText: /^Client secret$/}).locator('input[type="password"]').fill('grafana-secret');
    await page.locator('div').filter({hasText: /^Impersonation user$/}).locator('input').fill('service-account-grafana-client');
    await sel(page, 'Data source settings page Save and Test button').click();
    await page.waitForSelector('[role="alert"]', { timeout: 10000 });
}

async function runQueryAndCheckResults(page: Page) {
    if (isNewGrafana) {
        await page.getByLabel('Explore data').click();
    } else {
        await page.getByText('Explore', {exact: true}).click();
    }
    await page.getByTestId('data-testid TimePicker Open Button').click();
    await sel(page, 'Time Range from field').fill('1995-01-01');
    await sel(page, 'Time Range to field').fill('1995-12-31');
    await page.getByTestId('data-testid TimePicker submit button').click();
    if (isNewGrafana) {
        await page.locator('div').filter({hasText: /^Format asChoose$/}).locator('svg').click();
        await page.getByRole('option', {name: 'Table'}).click();
        await page.getByTestId('data-testid Code editor container').click();
    } else {
        await page.getByRole('combobox', { name: 'Format as' }).click();
        await page.getByRole('option', {name: 'Table'}).click();
    }
    await page.getByTestId('data-testid RefreshPicker run button').click();
    await expect(page.getByTestId('data-testid table body')).toContainText(/.*1995-01-19 0.:00:00.*/, { timeout: 30000 });
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
    // Alert is already visible (setupDataSourceWithClientCredentials waits for it).
    // Verify it is NOT the success message.
    await expect(page.locator('[role="alert"]').first()).toBeVisible({ timeout: 10000 });
    await expect(page.locator('[role="alert"]').filter({ hasText: 'Data source is working' })).toHaveCount(0);
});

test('test client credentials flow with configured access token', async ({ page }) => {
    await login(page);
    await goToTrinoSettings(page);
    await page.locator('div').filter({hasText: /^Access token$/}).locator('input[type="password"]').fill('aaa');
    await setupDataSourceWithClientCredentials(page, GRAFANA_CLIENT);
    // Alert is already visible (setupDataSourceWithClientCredentials waits for it).
    // Verify it is NOT the success message.
    await expect(page.locator('[role="alert"]').first()).toBeVisible({ timeout: 10000 });
    await expect(page.locator('[role="alert"]').filter({ hasText: 'Data source is working' })).toHaveCount(0);
});
