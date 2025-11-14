<script lang="ts">
	import Logo from '$components/icons/Logo.svelte';
	import SearchBar from '$components/SearchBar.svelte';
	import { fly } from 'svelte/transition';
	import { onMount } from 'svelte';
	import { themeManager } from '$lib/stores/theme.svelte';
	import { localeManager } from '$lib/stores/locale.svelte';
	import * as m from '$lib/paraglide/messages';
	import type { Locale } from '$lib/paraglide/runtime';
	import { page } from '$app/state';

	let showLocaleDropdown = $state(false);
	let isHomePage = $derived(page.url.pathname === '/');
	let searchQuery = $derived(page.url.searchParams.get('q') || '');

	function toggleLocaleDropdown() {
		showLocaleDropdown = !showLocaleDropdown;
	}

	async function selectLocale(locale: string) {
		showLocaleDropdown = false;
		await localeManager.switchLocale(locale as Locale);
	}

	// 点击外部关闭下拉框
	function handleClickOutside(event: MouseEvent) {
		const target = event.target as HTMLElement;
		if (!target.closest('.locale-dropdown')) {
			showLocaleDropdown = false;
		}
	}

	let themeToggleEvent: MouseEvent | undefined;
	// 捕获鼠标/指针点击位置用于动画，但不阻止 label 的默认行为
	function captureMousePosition(event: MouseEvent) {
		themeToggleEvent = event;
	}

	function handleThemeChange() {
		themeManager.setDark(themeManager.isDarkMode, themeToggleEvent);
		themeToggleEvent = undefined;
	}

	onMount(() => {
		localeManager.init();
		return themeManager.watchSystemTheme();
	});
</script>

<svelte:window onclick={handleClickOutside} />

<header class="sticky top-0 z-50 bg-base-100 shadow-sm">
	<div class="navbar mx-auto max-w-7xl px-3 sm:px-4 lg:px-6">
		<div class="navbar-start">
			{#if isHomePage}
				<div class="flex items-center gap-2 sm:gap-3">
					<div class="avatar">
						<div class="size-9 rounded-lg bg-primary sm:size-12 sm:rounded-xl">
							<Logo class="size-full p-1.5 text-primary-content sm:p-2" />
						</div>
					</div>
					<h1 class="text-lg font-bold sm:text-xl lg:text-2xl">{m.app_title()}</h1>
				</div>
			{:else}
				<a href="/" class="btn btn-circle btn-ghost btn-sm sm:btn-md" aria-label={m.back_home()}>
					<svg
						class="size-5 sm:size-6"
						fill="none"
						stroke="currentColor"
						viewBox="0 0 24 24"
						xmlns="http://www.w3.org/2000/svg"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M10 19l-7-7m0 0l7-7m-7 7h18"
						/>
					</svg>
				</a>
			{/if}
		</div>

		<div class="navbar-end gap-1 sm:gap-2">
			{#if isHomePage}
				<SearchBar initialQuery={searchQuery} />
			{/if}

			<!-- Language Switcher -->
			<div class="dropdown dropdown-end locale-dropdown">
				<button
					class="btn btn-circle btn-ghost btn-sm sm:btn-md"
					onclick={toggleLocaleDropdown}
					aria-label="Switch language"
				>
					{#await import('$components/icons/Globe.svelte') then { default: GlobeIcon }}
						<GlobeIcon class="size-5 sm:size-6" />
					{/await}
				</button>
				{#if showLocaleDropdown}
					<ul
						class="menu dropdown-content z-1 mt-3 w-40 rounded-box bg-base-100 p-2 shadow"
						transition:fly={{ y: -10, duration: 200 }}
					>
						{#each localeManager.availableLocales as locale}
							<li>
								<button
									onclick={() => selectLocale(locale.code)}
									class="text-sm"
									class:active={localeManager.currentLocale === locale.code}
								>
									<span>{locale.nativeName}</span>
									{#if localeManager.currentLocale === locale.code}
										<span class="ml-auto text-primary">✓</span>
									{/if}
								</button>
							</li>
						{/each}
					</ul>
				{/if}
			</div>

			<!-- Theme Switcher -->
			<label
				class="btn swap btn-circle swap-rotate btn-ghost btn-sm sm:btn-md"
				onpointerdown={captureMousePosition}
				aria-label={m.theme_toggle()}
			>
				<input
					type="checkbox"
					class="theme-controller hidden"
					bind:checked={themeManager.isDarkMode}
					onchange={handleThemeChange}
				/>
				{#await import('$components/icons/Sun.svelte') then { default: Sun }}
					<Sun class="swap-off size-5 sm:size-6" />
				{/await}
				{#await import('$components/icons/Moon.svelte') then { default: MoonIcon }}
					<MoonIcon class="swap-on size-5 sm:size-6" />
				{/await}
			</label>
		</div>
	</div>
</header>
