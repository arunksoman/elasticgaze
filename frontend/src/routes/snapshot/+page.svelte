<script>
	// Snapshot page component
</script>

<div class="p-6">
	<h1 class="text-2xl font-medium mb-4">Snapshots</h1>
	<div class="bg-fluent-surface-light dark:bg-fluent-surface-dark p-6 shadow-fluent">
		<p class="mb-4">View and manage your Elasticsearch snapshots</p>
		
		<div class="mb-4">
			<button class="fluent-btn bg-fluent-purple hover:bg-fluent-purple-light text-white px-4 py-2 mb-4">
				Create Snapshot
			</button>
		</div>
		
		<div class="border border-fluent-border-light dark:border-fluent-border-dark p-4 mb-4">
			<h2 class="text-lg font-medium mb-2">Snapshot Repositories</h2>
			<table class="w-full">
				<thead>
					<tr class="border-b border-fluent-border-light dark:border-fluent-border-dark">
						<th class="text-left p-2">Name</th>
						<th class="text-left p-2">Type</th>
						<th class="text-left p-2">Location</th>
					</tr>
				</thead>
				<tbody>
					<tr class="border-b border-fluent-border-light dark:border-fluent-border-dark">
						<td class="p-2">backup-repo</td>
						<td class="p-2">fs</td>
						<td class="p-2">/mnt/backups</td>
					</tr>
				</tbody>
			</table>
		</div>
		
		<div class="border border-fluent-border-light dark:border-fluent-border-dark p-4">
			<h2 class="text-lg font-medium mb-2">Snapshots</h2>
			<table class="w-full">
				<thead>
					<tr class="border-b border-fluent-border-light dark:border-fluent-border-dark">
						<th class="text-left p-2">Name</th>
						<th class="text-left p-2">Repository</th>
						<th class="text-left p-2">State</th>
						<th class="text-left p-2">Start Time</th>
					</tr>
				</thead>
				<tbody>
					<tr class="border-b border-fluent-border-light dark:border-fluent-border-dark">
						<td class="p-2">snapshot_1</td>
						<td class="p-2">backup-repo</td>
						<td class="p-2">SUCCESS</td>
						<td class="p-2">2025-09-26 14:30:00</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</div>