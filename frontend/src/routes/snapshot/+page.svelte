<script>
	// Snapshot page component
</script>

<div class="p-6">
	<h1 class="text-2xl font-medium mb-4 theme-text-primary">Snapshots</h1>
	<p class="mb-4 theme-text-secondary">View and manage your Elasticsearch snapshots</p>
	
	<div class="mb-4">
		<button class="fluent-btn bg-fluent-purple hover:bg-fluent-purple-light text-white px-4 py-2 mb-4 rounded">
			Create Snapshot
		</button>
	</div>
	
	<div class="border theme-border p-4 mb-4 rounded">
		<h2 class="text-lg font-medium mb-2 theme-text-primary">Snapshot Repositories</h2>
		<table class="w-full">
			<thead>
				<tr class="border-b theme-border">
					<th class="text-left p-2 theme-text-primary font-medium">Name</th>
					<th class="text-left p-2 theme-text-primary font-medium">Type</th>
					<th class="text-left p-2 theme-text-primary font-medium">Location</th>
				</tr>
			</thead>
			<tbody>
				<tr class="border-b theme-border">
					<td class="p-2 theme-text-secondary">backup-repo</td>
					<td class="p-2 theme-text-secondary">fs</td>
					<td class="p-2 theme-text-secondary">/mnt/backups</td>
				</tr>
			</tbody>
		</table>
	</div>
	
	<div class="border theme-border p-4 rounded">
		<h2 class="text-lg font-medium mb-2 theme-text-primary">Snapshots</h2>
		<table class="w-full">
			<thead>
				<tr class="border-b theme-border">
					<th class="text-left p-2 theme-text-primary font-medium">Name</th>
					<th class="text-left p-2 theme-text-primary font-medium">Repository</th>
					<th class="text-left p-2 theme-text-primary font-medium">State</th>
					<th class="text-left p-2 theme-text-primary font-medium">Start Time</th>
				</tr>
			</thead>
			<tbody>
				<tr class="border-b theme-border">
					<td class="p-2 theme-text-secondary">snapshot_1</td>
					<td class="p-2 theme-text-secondary">backup-repo</td>
					<td class="p-2 theme-text-secondary">SUCCESS</td>
					<td class="p-2 theme-text-secondary">2025-09-26 14:30:00</td>
				</tr>
			</tbody>
		</table>
	</div>
</div>