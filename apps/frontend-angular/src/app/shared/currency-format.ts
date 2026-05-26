import { CurrencyMetadata } from '../core/services/api.models';

export const defaultCurrency: CurrencyMetadata = {
  code: 'COP',
  symbol: '$',
  locale: 'es-CO',
  decimal_digits: 0,
  thousand_separator: '.',
  decimal_separator: ',',
  symbol_position: 'before',
};

export function formatClinicCurrency(value: number | null | undefined, currency: CurrencyMetadata | null | undefined): string {
  const amount = Number(value ?? 0);
  const metadata = currency ?? defaultCurrency;
  const formatted = new Intl.NumberFormat(metadata.locale, {
    style: 'currency',
    currency: metadata.code,
    currencyDisplay: 'symbol',
    minimumFractionDigits: metadata.decimal_digits,
    maximumFractionDigits: metadata.decimal_digits,
  }).format(amount);

  return `${metadata.code} ${formatted}`;
}
