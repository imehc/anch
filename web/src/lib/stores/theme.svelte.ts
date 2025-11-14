import { browser } from '$app/environment';
import { PUBLIC_LIGHT_THEME, PUBLIC_DARK_THEME } from '$env/static/public';

const LIGHT_THEME = PUBLIC_LIGHT_THEME;
const DARK_THEME = PUBLIC_DARK_THEME;

function getSystemTheme(): boolean {
	if (browser) {
		return window.matchMedia('(prefers-color-scheme: dark)').matches;
	}
	return false;
}

function getSavedTheme(): boolean {
	if (browser) {
		const saved = localStorage.getItem('theme');
		if (saved !== null) {
			return saved === DARK_THEME;
		}
		return getSystemTheme();
	}
	return false;
}

function applyTheme(dark: boolean, event?: MouseEvent) {
	if (browser) {
		const theme = dark ? DARK_THEME : LIGHT_THEME;

		// 检查浏览器是否支持 View Transitions API
		if (!document.startViewTransition) {
			document.documentElement.setAttribute('data-theme', theme);
			localStorage.setItem('theme', theme);
			return;
		}

		if (!event) {
			document.documentElement.setAttribute('data-theme', theme);
			localStorage.setItem('theme', theme);
			return;
		}

		// 获取点击位置
		const x = event.clientX;
		const y = event.clientY;

		// 计算从点击位置到页面最远角的距离
		const endRadius = Math.hypot(
			Math.max(x, innerWidth - x),
			Math.max(y, innerHeight - y)
		);

		// 设置 CSS 变量供动画使用
		document.documentElement.style.setProperty('--x', `${x}px`);
		document.documentElement.style.setProperty('--y', `${y}px`);
		document.documentElement.style.setProperty('--r', `${endRadius}px`);

		// 标记动画方向：切换到暗色主题用扩散，切换到亮色主题用收缩
		document.documentElement.setAttribute('data-theme-transition', dark ? 'dark' : 'light');

		// 使用 View Transitions API 实现圆形扩散动画
		const transition = document.startViewTransition(() => {
			document.documentElement.setAttribute('data-theme', theme);
			localStorage.setItem('theme', theme);
		});

		// 动画结束后清理标记
		transition.finished.finally(() => {
			document.documentElement.removeAttribute('data-theme-transition');
		});
	}
}

class ThemeManager {
	isDarkMode = $state(false);
	isInitialized = $state(false);

	init() {
		if (browser && !this.isInitialized) {
			this.isDarkMode = getSavedTheme();
			applyTheme(this.isDarkMode);
			this.isInitialized = true;
		}
	}

	toggle(event?: MouseEvent) {
		this.isDarkMode = !this.isDarkMode;
		applyTheme(this.isDarkMode, event);
	}

	setDark(dark: boolean, event?: MouseEvent) {
		this.isDarkMode = dark;
		applyTheme(this.isDarkMode, event);
	}

	watchSystemTheme() {
		this.init()
		if (browser) {
			const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
			const handleChange = (e: MediaQueryListEvent) => {
				if (localStorage.getItem('theme') === null) {
					this.setDark(e.matches);
				}
			};
			mediaQuery.addEventListener('change', handleChange);
			return () => mediaQuery.removeEventListener('change', handleChange);
		}
		return () => { };
	}
}

export const themeManager = new ThemeManager();
export { LIGHT_THEME, DARK_THEME };
