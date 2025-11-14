<script lang="ts">
	import CopyIcon from '$components/icons/Copy.svelte';
	import CheckIcon from '$components/icons/Check.svelte';
	import * as m from '$lib/paraglide/messages';
	import {
		hexToRgba,
		hslToRgb,
		hsvToRgb,
		oklchToRgb,
		rgbaToString,
		rgbaToHex,
		rgbToHex,
		hslaToString,
		hsvaToString,
		oklchaToString,
		isValidHex,
		parseRgbaString,
		parseHslaString,
		parseHsvaString,
		parseOklchaString,
		rgbToRgba,
		rgbToHsl,
		rgbToHsv,
		rgbToOklch,
		hslToHsla,
		hsvToHsva,
		oklchToOklcha,
		type RGB,
		type RGBA
	} from '$lib/utils/color';

	// 主颜色状态（RGB）
	let rgb = $state<RGB>({ r: 66, g: 42, b: 213 });
	let alpha = $state(1); // 透明度 0-1
	let copySuccess = $state('');
	let copyTimeout: ReturnType<typeof setTimeout> | undefined;

	// 输入框状态
	let hexInput = $state('#422ad5ff');
	let rgbaInput = $state('');
	let hslaInput = $state('');
	let hsvaInput = $state('');
	let oklchaInput = $state('');

	// 计算派生值
	let hsl = $derived(rgbToHsl(rgb));
	let hsv = $derived(rgbToHsv(rgb));
	let oklch = $derived(rgbToOklch(rgb));
	let rgba = $derived<RGBA>(rgbToRgba(rgb, alpha));
	let hsla = $derived(hslToHsla(hsl, alpha));
	let hsva = $derived(hsvToHsva(hsv, alpha));
	let oklcha = $derived(oklchToOklcha(oklch, alpha));

	let hexString = $derived(rgbaToHex(rgba));
	let rgbaString = $derived(rgbaToString(rgba));
	let hslaString = $derived(hslaToString(hsla));
	let hsvaString = $derived(hsvaToString(hsva));
	let oklchaString = $derived(oklchaToString(oklcha));

	// 同步输入框
	$effect(() => {
		hexInput = hexString;
		rgbaInput = rgbaString;
		hslaInput = hslaString;
		hsvaInput = hsvaString;
		oklchaInput = oklchaString;
	});

	// HEX 输入处理
	function handleHexInput(e: Event) {
		const input = e.target as HTMLInputElement;
		const value = input.value;
		if (value.match(/^#?([0-9A-Fa-f]{0,8})$/)) {
			if (isValidHex(value)) {
				const parsed = hexToRgba(value);
				rgb = { r: parsed.r, g: parsed.g, b: parsed.b };
				alpha = parsed.a;
			}
		}
	}

	// RGBA 输入处理
	function handleRgbaInput(e: Event) {
		const input = e.target as HTMLInputElement;
		const value = input.value;
		const parsed = parseRgbaString(value);
		if (parsed) {
			rgb = { r: parsed.r, g: parsed.g, b: parsed.b };
			alpha = parsed.a;
		}
	}

	// HSLA 输入处理
	function handleHslaInput(e: Event) {
		const input = e.target as HTMLInputElement;
		const value = input.value;
		const parsed = parseHslaString(value);
		if (parsed) {
			rgb = hslToRgb({ h: parsed.h, s: parsed.s, l: parsed.l });
			alpha = parsed.a;
		}
	}

	// HSVA 输入处理
	function handleHsvaInput(e: Event) {
		const input = e.target as HTMLInputElement;
		const value = input.value;
		const parsed = parseHsvaString(value);
		if (parsed) {
			rgb = hsvToRgb({ h: parsed.h, s: parsed.s, v: parsed.v });
			alpha = parsed.a;
		}
	}

	// OKLCHA 输入处理
	function handleOklchaInput(e: Event) {
		const input = e.target as HTMLInputElement;
		const value = input.value;
		const parsed = parseOklchaString(value);
		if (parsed) {
			try {
				rgb = oklchToRgb({ l: parsed.l, c: parsed.c, h: parsed.h });
				alpha = parsed.a;
			} catch (err) {
				console.error('Invalid OKLCHA color:', err);
			}
		}
	}

	// 复制到剪贴板
	async function copyToClipboard(text: string, type: string) {
		try {
			await navigator.clipboard.writeText(text);
			copySuccess = type;
			if (copyTimeout) clearTimeout(copyTimeout);
			copyTimeout = setTimeout(() => {
				copySuccess = '';
			}, 2000);
		} catch (err) {
			console.error('Failed to copy:', err);
		}
	}

	// 预设颜色
	const presetColors = [
		'#FF0000', // Red
		'#FF7F00', // Orange
		'#FFFF00', // Yellow
		'#00FF00', // Green
		'#0000FF', // Blue
		'#4B0082', // Indigo
		'#9400D3', // Violet
		'#FF1493', // Deep Pink
		'#00CED1', // Dark Turquoise
		'#FFD700', // Gold
		'#FF69B4', // Hot Pink
		'#8B4513', // Saddle Brown
		'#000000', // Black
		'#808080', // Gray
		'#FFFFFF' // White
	];

	let colorPickerInput: HTMLInputElement | undefined;

	function applyPreset(color: string) {
		const parsed = hexToRgba(color + 'FF'); // Add FF for full opacity
		rgb = { r: parsed.r, g: parsed.g, b: parsed.b };
		alpha = 1;
	}

	function handleColorPickerChange(e: Event) {
		const input = e.target as HTMLInputElement;
		const value = input.value;
		const parsed = hexToRgba(value + 'FF');
		rgb = { r: parsed.r, g: parsed.g, b: parsed.b };
		alpha = 1;
	}

	function openColorPicker() {
		colorPickerInput?.click();
	}
</script>

{#snippet colorInput(id: string, label: string, inputValue: string, placeholder: string, readonly: boolean, inputHandler: ((e: Event) => void) | undefined, copyValue: string, copyType: string, ariaLabel: string)}
	<div>
		<label class="label" for={id}>
			<span class="label-text font-semibold">{label}</span>
		</label>
		<div class="join w-full">
			<input
				{id}
				type="text"
				{readonly}
				value={inputValue}
				oninput={inputHandler}
				class="input join-item input-bordered w-full font-mono text-sm {id === 'hex-input'
					? 'uppercase'
					: ''}"
				{placeholder}
			/>
			<button
				onclick={() => copyToClipboard(copyValue, copyType)}
				class="btn join-item btn-primary"
				class:btn-success={copySuccess === copyType}
				aria-label={ariaLabel}
			>
				{#if copySuccess === copyType}
					<CheckIcon class="size-5" />
				{:else}
					<CopyIcon class="size-5" />
				{/if}
			</button>
		</div>
	</div>
{/snippet}

<!-- Main Content -->
<main class="mx-auto max-w-4xl px-3 py-4 sm:px-4 sm:py-6 lg:px-6 lg:py-8">
	<!-- Page Title -->
	<div class="mb-6 sm:mb-8">
		<h1 class="mb-2 text-2xl font-bold sm:text-3xl lg:text-4xl">
			{m.tool_color_title()}
		</h1>
		<p class="text-sm text-base-content/70 sm:text-base">
			{m.tool_color_desc()}
		</p>
	</div>

	<!-- Color Display Card -->
	<div class="card mb-6 bg-base-100 shadow-xl">
		<div class="card-body p-4 sm:p-6">
			<div class="mb-6">
				<div
					class="h-48 w-full rounded-lg shadow-inner transition-colors duration-150 sm:h-64"
					style="background-color: {hexString};"
				></div>
			</div>

			<!-- Color Values Display -->
			<div class="mb-6 grid gap-4 sm:grid-cols-2">
				{@render colorInput('rgba-input', 'RGBA', rgbaInput, 'rgba(0, 0, 0, 1)', false, handleRgbaInput, rgbaString, 'rgba', 'Copy RGBA')}
				{@render colorInput('hex-input', 'HEX (8-digit)', hexInput, '#00000000', false, handleHexInput, hexString, 'hex', m.color_copy_hex())}
				{@render colorInput('hsla-input', 'HSLA', hslaInput, 'hsla(0, 100%, 50%, 1)', false, handleHslaInput, hslaString, 'hsla', 'Copy HSLA')}
				{@render colorInput('hsva-input', 'HSVA', hsvaInput, 'hsva(0, 100%, 100%, 1)', false, handleHsvaInput, hsvaString, 'hsva', 'Copy HSVA')}
				<div class="sm:col-span-2">
					{@render colorInput('oklcha-input', 'OKLCH + Alpha', oklchaInput, 'oklch(0.5 0.2 180 / 1)', false, handleOklchaInput, oklchaString, 'oklcha', 'Copy OKLCH+A')}
				</div>
			</div>

			<!-- RGB Sliders -->
			<div class="space-y-4">
				<!-- Red Slider -->
				<div>
					<div class="mb-2 flex items-center justify-between">
						<label for="red-slider" class="label-text font-semibold text-red-600">Red</label>
						<span class="text-sm font-mono">{rgb.r}</span>
					</div>
					<input
						id="red-slider"
						type="range"
						min="0"
						max="255"
						bind:value={rgb.r}
						class="range range-error w-full"
					/>
				</div>

				<!-- Green Slider -->
				<div>
					<div class="mb-2 flex items-center justify-between">
						<label for="green-slider" class="label-text font-semibold text-green-600"
							>Green</label
						>
						<span class="text-sm font-mono">{rgb.g}</span>
					</div>
					<input
						id="green-slider"
						type="range"
						min="0"
						max="255"
						bind:value={rgb.g}
						class="range range-success w-full"
					/>
				</div>

				<!-- Blue Slider -->
				<div>
					<div class="mb-2 flex items-center justify-between">
						<label for="blue-slider" class="label-text font-semibold text-blue-600">Blue</label>
						<span class="text-sm font-mono">{rgb.b}</span>
					</div>
					<input
						id="blue-slider"
						type="range"
						min="0"
						max="255"
						bind:value={rgb.b}
						class="range range-info w-full"
					/>
				</div>

				<!-- Alpha Slider -->
				<div>
					<div class="mb-2 flex items-center justify-between">
						<label for="alpha-slider" class="label-text font-semibold">Alpha</label>
						<span class="text-sm font-mono">{Math.round(alpha * 100) / 100}</span>
					</div>
					<input
						id="alpha-slider"
						type="range"
						min="0"
						max="1"
						step="0.01"
						bind:value={alpha}
						class="range w-full"
					/>
				</div>
			</div>
		</div>
		</div>

		<!-- Preset Colors -->
		<div class="card bg-base-100 shadow-xl">
			<div class="card-body p-4 sm:p-6">
				<h2 class="card-title mb-4 text-lg sm:text-xl">{m.color_preset_title()}</h2>
				<div class="grid grid-cols-5 gap-2 sm:grid-cols-8 sm:gap-3">
					{#each presetColors as color}
						<button
							onclick={() => applyPreset(color)}
							class="aspect-square rounded-lg border-2 border-base-300 shadow-sm transition-transform hover:scale-110 hover:shadow-md active:scale-95"
							style="background-color: {color};"
							aria-label="Apply color {color}"
						></button>
					{/each}
					<!-- Custom Color Picker -->
					<button
						onclick={openColorPicker}
						class="aspect-square rounded-lg border-2 border-dashed border-base-300 shadow-sm transition-transform hover:scale-110 hover:shadow-md active:scale-95 flex items-center justify-center text-base-content/50 hover:text-base-content/70"
						aria-label="Choose custom color"
					>
						<svg class="size-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
						</svg>
					</button>
					<input
						bind:this={colorPickerInput}
						type="color"
						onchange={handleColorPickerChange}
						class="hidden"
						value={rgbToHex(rgb)}
					/>
				</div>
			</div>
		</div>
</main>
