import { browser } from '$app/environment';
import { locales, baseLocale, setLocale, getLocale, type Locale } from '$lib/paraglide/runtime';

interface LocaleInfo {
	code: Locale;
	name: string;
	nativeName: string;
}

const LOCALE_NAMES: Record<Locale, { name: string; nativeName: string }> = {
	en: { name: 'English', nativeName: 'English' },
	zh: { name: 'Chinese', nativeName: '简体中文' }
};

function getAvailableLocales(): LocaleInfo[] {
	return locales.map((code) => ({
		code,
		name: LOCALE_NAMES[code].name,
		nativeName: LOCALE_NAMES[code].nativeName
	}));
}

class LocaleManager {
	currentLocale = $state<Locale>(baseLocale);
	availableLocales = $state<LocaleInfo[]>(getAvailableLocales());
	isInitialized = $state(false);

	init() {
		if (browser && !this.isInitialized) {
			this.currentLocale = getLocale();
			this.isInitialized = true;
		}
	}

	async switchLocale(newLocale: Locale) {
		if (newLocale === this.currentLocale) {
			return;
		}

		this.currentLocale = newLocale;

		if (browser) {
			// setLocale will handle cookie/URL updates and page reload
			await setLocale(newLocale);
		}
	}

	getLocaleName(code: Locale): string {
		return LOCALE_NAMES[code]?.name || code;
	}

	getNativeLocaleName(code: Locale): string {
		return LOCALE_NAMES[code]?.nativeName || code;
	}
}

export const localeManager = new LocaleManager();
