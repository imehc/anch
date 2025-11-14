<script lang="ts">
	import SearchIcon from '$components/icons/Search.svelte';
	import { fly } from 'svelte/transition';
	import { cubicOut } from 'svelte/easing';
	import { goto } from '$app/navigation';
	import * as m from '$lib/paraglide/messages';

	interface SearchBarProps {
		initialQuery?: string;
	}

	let { initialQuery = '' }: SearchBarProps = $props();

	let showMobileSearch = $state(initialQuery !== '');
	let searchQuery = $state(initialQuery);
	let searchInput = $state<HTMLInputElement>();

	// 只在 initialQuery 变化时同步（不影响用户输入）
	let lastInitialQuery = $state(initialQuery);
	$effect(() => {
		if (initialQuery !== lastInitialQuery) {
			searchQuery = initialQuery;
			showMobileSearch = initialQuery !== '';
			lastInitialQuery = initialQuery;
		}
	});

	function toggleMobileSearch() {
		showMobileSearch = !showMobileSearch;
		if (!showMobileSearch) {
			searchQuery = '';
			goto('/');
		} else {
			// 延迟聚焦，等待动画完成
			setTimeout(() => {
				searchInput?.focus();
			}, 100);
		}
	}

	function handleSearch(e: Event) {
		e.preventDefault();
		if (searchQuery.trim()) {
			goto(`/?q=${encodeURIComponent(searchQuery.trim())}`);
		} else {
			goto('/');
		}
	}
</script>

{#if showMobileSearch}
	<!-- Mobile Search Bar (animated) -->
	<form
		onsubmit={handleSearch}
		class="flex-1 origin-right lg:hidden"
		in:fly={{ x: 50, duration: 200, easing: cubicOut }}
		out:fly={{ x: 50, duration: 150, easing: cubicOut }}
	>
		<div class="relative">
			<input
				bind:this={searchInput}
				type="text"
				bind:value={searchQuery}
				placeholder={m.search_placeholder_mobile()}
				class="input-bordered input input-sm w-full pr-10 pl-9 sm:input-md sm:pr-12 sm:pl-10"
			/>
			<span class="pointer-events-none absolute top-1/2 left-2.5 z-10 -translate-y-1/2 sm:left-3">
				<SearchIcon class="size-4 text-base-content/40 sm:size-5" />
			</span>
			<button
				type="button"
				onclick={toggleMobileSearch}
				class="btn absolute top-1/2 right-1.5 z-10 btn-circle -translate-y-1/2 btn-ghost btn-xs sm:right-2 sm:btn-sm"
				aria-label={m.close_search()}
			>
				{#await import('$components/icons/Close.svelte') then { default: ColseIcon }}
					<ColseIcon class="size-4 sm:size-5" />
				{/await}
			</button>
		</div>
	</form>
{:else}
	<!-- Mobile Search Button -->
	<button
		class="btn btn-circle btn-ghost btn-sm sm:btn-md lg:hidden"
		aria-label={m.search_label()}
		onclick={toggleMobileSearch}
	>
		<SearchIcon class="size-5 sm:size-6" />
	</button>
{/if}

<!-- PC Search Bar -->
<form onsubmit={handleSearch} class="hidden lg:block">
	<div class="relative">
		<input
			type="text"
			bind:value={searchQuery}
			placeholder={m.search_placeholder()}
			class="input-bordered input input-md w-64 pr-4 pl-10 transition-all duration-300 focus:w-80 focus:shadow-md xl:w-80"
		/>
		<span class="pointer-events-none absolute top-1/2 left-3 z-10 -translate-y-1/2">
			<SearchIcon class="size-5 text-base-content/40" />
		</span>
	</div>
</form>
