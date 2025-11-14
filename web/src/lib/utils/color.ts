/**
 * Color conversion utilities
 * Supports RGB, HEX, HSL, HSV, and OKLCH formats
 */

export interface RGB {
	r: number; // 0-255
	g: number; // 0-255
	b: number; // 0-255
}

export interface RGBA {
	r: number; // 0-255
	g: number; // 0-255
	b: number; // 0-255
	a: number; // 0-1
}

export interface HSL {
	h: number; // 0-360
	s: number; // 0-100
	l: number; // 0-100
}

export interface HSLA {
	h: number; // 0-360
	s: number; // 0-100
	l: number; // 0-100
	a: number; // 0-1
}

export interface HSV {
	h: number; // 0-360
	s: number; // 0-100
	v: number; // 0-100
}

export interface HSVA {
	h: number; // 0-360
	s: number; // 0-100
	v: number; // 0-100
	a: number; // 0-1
}

export interface OKLCH {
	l: number; // 0-1
	c: number; // 0-0.4 (typically)
	h: number; // 0-360
}

export interface OKLCHA {
	l: number; // 0-1
	c: number; // 0-0.4 (typically)
	h: number; // 0-360
	a: number; // 0-1
}

// RGB to HEX
export function rgbToHex(rgb: RGB): string {
	const toHex = (n: number) => {
		const hex = Math.round(n).toString(16);
		return hex.length === 1 ? '0' + hex : hex;
	};
	return '#' + toHex(rgb.r) + toHex(rgb.g) + toHex(rgb.b);
}

// RGBA to 8-digit HEX
export function rgbaToHex(rgba: RGBA): string {
	const toHex = (n: number) => {
		const hex = Math.round(n).toString(16);
		return hex.length === 1 ? '0' + hex : hex;
	};
	const alphaHex = toHex(Math.round(rgba.a * 255));
	return '#' + toHex(rgba.r) + toHex(rgba.g) + toHex(rgba.b) + alphaHex;
}

// HEX to RGB
export function hexToRgb(hex: string): RGB {
	hex = hex.replace('#', '');
	if (hex.length === 3) {
		hex = hex
			.split('')
			.map((char) => char + char)
			.join('');
	}
	return {
		r: parseInt(hex.substring(0, 2), 16),
		g: parseInt(hex.substring(2, 4), 16),
		b: parseInt(hex.substring(4, 6), 16)
	};
}

// HEX to RGBA (supports 8-digit HEX)
export function hexToRgba(hex: string): RGBA {
	hex = hex.replace('#', '');
	if (hex.length === 3) {
		hex = hex
			.split('')
			.map((char) => char + char)
			.join('');
	}

	const r = parseInt(hex.substring(0, 2), 16);
	const g = parseInt(hex.substring(2, 4), 16);
	const b = parseInt(hex.substring(4, 6), 16);
	const a = hex.length === 8 ? parseInt(hex.substring(6, 8), 16) / 255 : 1;

	return { r, g, b, a };
}

// RGB to HSL
export function rgbToHsl(rgb: RGB): HSL {
	const r = rgb.r / 255;
	const g = rgb.g / 255;
	const b = rgb.b / 255;

	const max = Math.max(r, g, b);
	const min = Math.min(r, g, b);
	const delta = max - min;

	let h = 0;
	let s = 0;
	const l = (max + min) / 2;

	if (delta !== 0) {
		s = l > 0.5 ? delta / (2 - max - min) : delta / (max + min);

		switch (max) {
			case r:
				h = ((g - b) / delta + (g < b ? 6 : 0)) / 6;
				break;
			case g:
				h = ((b - r) / delta + 2) / 6;
				break;
			case b:
				h = ((r - g) / delta + 4) / 6;
				break;
		}
	}

	return {
		h: Math.round(h * 360),
		s: Math.round(s * 100),
		l: Math.round(l * 100)
	};
}

// HSL to RGB
export function hslToRgb(hsl: HSL): RGB {
	const h = hsl.h / 360;
	const s = hsl.s / 100;
	const l = hsl.l / 100;

	let r: number, g: number, b: number;

	if (s === 0) {
		r = g = b = l;
	} else {
		const hue2rgb = (p: number, q: number, t: number) => {
			if (t < 0) t += 1;
			if (t > 1) t -= 1;
			if (t < 1 / 6) return p + (q - p) * 6 * t;
			if (t < 1 / 2) return q;
			if (t < 2 / 3) return p + (q - p) * (2 / 3 - t) * 6;
			return p;
		};

		const q = l < 0.5 ? l * (1 + s) : l + s - l * s;
		const p = 2 * l - q;

		r = hue2rgb(p, q, h + 1 / 3);
		g = hue2rgb(p, q, h);
		b = hue2rgb(p, q, h - 1 / 3);
	}

	return {
		r: Math.round(r * 255),
		g: Math.round(g * 255),
		b: Math.round(b * 255)
	};
}

// RGB to HSV
export function rgbToHsv(rgb: RGB): HSV {
	const r = rgb.r / 255;
	const g = rgb.g / 255;
	const b = rgb.b / 255;

	const max = Math.max(r, g, b);
	const min = Math.min(r, g, b);
	const delta = max - min;

	let h = 0;
	let s = max === 0 ? 0 : delta / max;
	const v = max;

	if (delta !== 0) {
		switch (max) {
			case r:
				h = ((g - b) / delta + (g < b ? 6 : 0)) / 6;
				break;
			case g:
				h = ((b - r) / delta + 2) / 6;
				break;
			case b:
				h = ((r - g) / delta + 4) / 6;
				break;
		}
	}

	return {
		h: Math.round(h * 360),
		s: Math.round(s * 100),
		v: Math.round(v * 100)
	};
}

// HSV to RGB
export function hsvToRgb(hsv: HSV): RGB {
	const h = hsv.h / 360;
	const s = hsv.s / 100;
	const v = hsv.v / 100;

	const i = Math.floor(h * 6);
	const f = h * 6 - i;
	const p = v * (1 - s);
	const q = v * (1 - f * s);
	const t = v * (1 - (1 - f) * s);

	let r: number, g: number, b: number;

	switch (i % 6) {
		case 0:
			r = v;
			g = t;
			b = p;
			break;
		case 1:
			r = q;
			g = v;
			b = p;
			break;
		case 2:
			r = p;
			g = v;
			b = t;
			break;
		case 3:
			r = p;
			g = q;
			b = v;
			break;
		case 4:
			r = t;
			g = p;
			b = v;
			break;
		case 5:
		default:
			r = v;
			g = p;
			b = q;
			break;
	}

	return {
		r: Math.round(r * 255),
		g: Math.round(g * 255),
		b: Math.round(b * 255)
	};
}

// Linear RGB helpers for OKLCH conversion
function srgbToLinear(c: number): number {
	const abs = Math.abs(c);
	if (abs <= 0.04045) {
		return c / 12.92;
	}
	return (Math.sign(c) * Math.pow((abs + 0.055) / 1.055, 2.4));
}

function linearToSrgb(c: number): number {
	const abs = Math.abs(c);
	if (abs > 0.0031308) {
		return Math.sign(c) * (1.055 * Math.pow(abs, 1 / 2.4) - 0.055);
	}
	return 12.92 * c;
}

// RGB to Linear RGB
function rgbToLinearRgb(rgb: RGB): [number, number, number] {
	return [srgbToLinear(rgb.r / 255), srgbToLinear(rgb.g / 255), srgbToLinear(rgb.b / 255)];
}

// Linear RGB to RGB
function linearRgbToRgb(linear: [number, number, number]): RGB {
	return {
		r: Math.round(Math.max(0, Math.min(255, linearToSrgb(linear[0]) * 255))),
		g: Math.round(Math.max(0, Math.min(255, linearToSrgb(linear[1]) * 255))),
		b: Math.round(Math.max(0, Math.min(255, linearToSrgb(linear[2]) * 255)))
	};
}

// Linear RGB to XYZ (D65)
function linearRgbToXyz(rgb: [number, number, number]): [number, number, number] {
	const [r, g, b] = rgb;
	return [
		0.4124564 * r + 0.3575761 * g + 0.1804375 * b,
		0.2126729 * r + 0.7151522 * g + 0.072175 * b,
		0.0193339 * r + 0.119192 * g + 0.9503041 * b
	];
}

// XYZ to Linear RGB (D65)
function xyzToLinearRgb(xyz: [number, number, number]): [number, number, number] {
	const [x, y, z] = xyz;
	return [
		3.2404542 * x - 1.5371385 * y - 0.4985314 * z,
		-0.969266 * x + 1.8760108 * y + 0.041556 * z,
		0.0556434 * x - 0.2040259 * y + 1.0572252 * z
	];
}

// XYZ to OKLab
function xyzToOklab(xyz: [number, number, number]): [number, number, number] {
	const [x, y, z] = xyz;

	const l = 0.8189330101 * x + 0.3618667424 * y - 0.1288597137 * z;
	const m = 0.0329845436 * x + 0.9293118715 * y + 0.0361456387 * z;
	const s = 0.0482003018 * x + 0.2643662691 * y + 0.6338517070 * z;

	const l_ = Math.cbrt(l);
	const m_ = Math.cbrt(m);
	const s_ = Math.cbrt(s);

	return [
		0.2104542553 * l_ + 0.793617785 * m_ - 0.0040720468 * s_,
		1.9779984951 * l_ - 2.428592205 * m_ + 0.4505937099 * s_,
		0.0259040371 * l_ + 0.7827717662 * m_ - 0.808675766 * s_
	];
}

// OKLab to XYZ
function oklabToXyz(lab: [number, number, number]): [number, number, number] {
	const [L, a, b] = lab;

	const l_ = L + 0.3963377774 * a + 0.2158037573 * b;
	const m_ = L - 0.1055613458 * a - 0.0638541728 * b;
	const s_ = L - 0.0894841775 * a - 1.291485548 * b;

	const l = l_ * l_ * l_;
	const m = m_ * m_ * m_;
	const s = s_ * s_ * s_;

	return [
		1.227013851103521026 * l - 0.5577999806518222383 * m + 0.28125614896646780758 * s,
		-0.040580178423280593977 * l + 1.1122568696168301049 * m - 0.071676678665601200577 * s,
		-0.076381284505706892869 * l - 0.4214819784180127136 * m + 1.5861632204407947575 * s
	];
}

// RGB to OKLCH
export function rgbToOklch(rgb: RGB): OKLCH {
	const linear = rgbToLinearRgb(rgb);
	const xyz = linearRgbToXyz(linear);
	const [L, a, b] = xyzToOklab(xyz);

	const c = Math.sqrt(a * a + b * b);
	let h = Math.atan2(b, a) * (180 / Math.PI);
	if (h < 0) h += 360;

	return {
		l: Math.round(L * 1000) / 1000,
		c: Math.round(c * 1000) / 1000,
		h: Math.round(h * 10) / 10
	};
}

// OKLCH to RGB
export function oklchToRgb(oklch: OKLCH): RGB {
	const { l, c, h } = oklch;

	const hRad = (h * Math.PI) / 180;
	const a = c * Math.cos(hRad);
	const b = c * Math.sin(hRad);

	const xyz = oklabToXyz([l, a, b]);
	const linear = xyzToLinearRgb(xyz);
	return linearRgbToRgb(linear);
}

// Format converters
export function rgbToString(rgb: RGB): string {
	return `rgb(${rgb.r}, ${rgb.g}, ${rgb.b})`;
}

export function rgbaToString(rgba: RGBA): string {
	return `rgba(${rgba.r}, ${rgba.g}, ${rgba.b}, ${Math.round(rgba.a * 100) / 100})`;
}

export function hslToString(hsl: HSL): string {
	return `hsl(${hsl.h}, ${hsl.s}%, ${hsl.l}%)`;
}

export function hslaToString(hsla: HSLA): string {
	return `hsla(${hsla.h}, ${hsla.s}%, ${hsla.l}%, ${Math.round(hsla.a * 100) / 100})`;
}

export function hsvToString(hsv: HSV): string {
	return `hsv(${hsv.h}, ${hsv.s}%, ${hsv.v}%)`;
}

export function hsvaToString(hsva: HSVA): string {
	return `hsva(${hsva.h}, ${hsva.s}%, ${hsva.v}%, ${Math.round(hsva.a * 100) / 100})`;
}

export function oklchToString(oklch: OKLCH): string {
	return `oklch(${oklch.l} ${oklch.c} ${oklch.h})`;
}

export function oklchaToString(oklcha: OKLCHA): string {
	return `oklch(${oklcha.l} ${oklcha.c} ${oklcha.h} / ${Math.round(oklcha.a * 100) / 100})`;
}

// Validate HEX color (supports 3, 6, and 8 digit HEX)
export function isValidHex(hex: string): boolean {
	hex = hex.replace('#', '');
	return /^[0-9A-Fa-f]{3}$|^[0-9A-Fa-f]{6}$|^[0-9A-Fa-f]{8}$/.test(hex);
}

// Parse RGBA string
export function parseRgbaString(rgba: string): RGBA | null {
	// Match: rgba(255, 0, 0, 0.5) or rgba(255 0 0 / 0.5)
	const match = rgba.match(
		/rgba?\((\d+),?\s*(\d+),?\s*(\d+)(?:,?\s*\/?\s*([\d.]+))?\)/i
	);
	if (match) {
		const r = parseInt(match[1]);
		const g = parseInt(match[2]);
		const b = parseInt(match[3]);
		const a = match[4] ? parseFloat(match[4]) : 1;

		if (
			r >= 0 &&
			r <= 255 &&
			g >= 0 &&
			g <= 255 &&
			b >= 0 &&
			b <= 255 &&
			a >= 0 &&
			a <= 1
		) {
			return { r, g, b, a };
		}
	}
	return null;
}

// Convert RGBA to RGB (removing alpha)
export function rgbaToRgb(rgba: RGBA): RGB {
	return { r: rgba.r, g: rgba.g, b: rgba.b };
}

// Convert RGB to RGBA (adding alpha)
export function rgbToRgba(rgb: RGB, alpha: number = 1): RGBA {
	return { r: rgb.r, g: rgb.g, b: rgb.b, a: alpha };
}

// Convert HSL to HSLA (adding alpha)
export function hslToHsla(hsl: HSL, alpha: number = 1): HSLA {
	return { h: hsl.h, s: hsl.s, l: hsl.l, a: alpha };
}

// Convert HSV to HSVA (adding alpha)
export function hsvToHsva(hsv: HSV, alpha: number = 1): HSVA {
	return { h: hsv.h, s: hsv.s, v: hsv.v, a: alpha };
}

// Convert OKLCH to OKLCHA (adding alpha)
export function oklchToOklcha(oklch: OKLCH, alpha: number = 1): OKLCHA {
	return { l: oklch.l, c: oklch.c, h: oklch.h, a: alpha };
}

// Parse HSLA string
export function parseHslaString(hsla: string): HSLA | null {
	const match = hsla.match(/hsla?\((\d+),?\s*(\d+)%?,?\s*(\d+)%?(?:,?\s*\/?\s*([\d.]+))?\)/i);
	if (match) {
		const h = parseInt(match[1]);
		const s = parseInt(match[2]);
		const l = parseInt(match[3]);
		const a = match[4] ? parseFloat(match[4]) : 1;

		if (h >= 0 && h <= 360 && s >= 0 && s <= 100 && l >= 0 && l <= 100 && a >= 0 && a <= 1) {
			return { h, s, l, a };
		}
	}
	return null;
}

// Parse HSVA string
export function parseHsvaString(hsva: string): HSVA | null {
	const match = hsva.match(/hsva?\((\d+),?\s*(\d+)%?,?\s*(\d+)%?(?:,?\s*\/?\s*([\d.]+))?\)/i);
	if (match) {
		const h = parseInt(match[1]);
		const s = parseInt(match[2]);
		const v = parseInt(match[3]);
		const a = match[4] ? parseFloat(match[4]) : 1;

		if (h >= 0 && h <= 360 && s >= 0 && s <= 100 && v >= 0 && v <= 100 && a >= 0 && a <= 1) {
			return { h, s, v, a };
		}
	}
	return null;
}

// Parse OKLCHA string
export function parseOklchaString(oklcha: string): OKLCHA | null {
	const match = oklcha.match(/oklch\(([\d.]+%?)\s+([\d.]+)\s+([\d.]+)(?:\s*\/\s*([\d.]+))?\)/i);
	if (match) {
		let l = parseFloat(match[1]);
		if (match[1].includes('%')) {
			l = l / 100;
		}
		const c = parseFloat(match[2]);
		const h = parseFloat(match[3]);
		const a = match[4] ? parseFloat(match[4]) : 1;

		if (l >= 0 && l <= 1 && c >= 0 && h >= 0 && h <= 360 && a >= 0 && a <= 1) {
			return { l, c, h, a };
		}
	}
	return null;
}
