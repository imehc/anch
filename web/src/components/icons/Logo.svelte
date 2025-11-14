<script lang="ts">
	let { class: className = 'size-6' }: { class?: string } = $props();
</script>

<svg class={className} viewBox="0 0 100 100" fill="none" xmlns="http://www.w3.org/2000/svg">
	<!-- 外圆形背景 -->
	<circle cx="50" cy="50" r="45" fill="currentColor" opacity="0.1" />

	<!-- TO 文字 -->
	<text
		x="50"
		y="42"
		font-family="system-ui, -apple-system, sans-serif"
		font-size="26"
		font-weight="700"
		fill="currentColor"
		text-anchor="middle"
		letter-spacing="-1"
	>
		TO
	</text>

	<!-- 分隔线 -->
	<line x1="25" y1="50" x2="75" y2="50" stroke="currentColor" stroke-width="2" opacity="0.3" />

	<!-- GO 文字 -->
	<text
		x="50"
		y="72"
		font-family="system-ui, -apple-system, sans-serif"
		font-size="26"
		font-weight="700"
		fill="currentColor"
		text-anchor="middle"
		letter-spacing="-1"
	>
		GO
	</text>
</svg>
