import elasticsearchMethodsPaths from '$lib/data/elasticsearch_methods_paths.json';

export interface UrlSuggestion {
	path: string;
	score: number;
}

export class UrlAutocompleteService {
	/**
	 * Get URL suggestions based on method and current input
	 */
	static getSuggestions(method: string, input: string, maxResults: number = 10): UrlSuggestion[] {
		const methodPaths = elasticsearchMethodsPaths[method as keyof typeof elasticsearchMethodsPaths] || [];
		
		if (!input.trim()) {
			// Return all paths for the method, sorted by length (shorter first)
			return methodPaths
				.map(path => ({ path, score: 1 }))
				.sort((a, b) => a.path.length - b.path.length)
				.slice(0, maxResults);
		}
		
		const normalizedInput = input.toLowerCase().trim();
		const suggestions: UrlSuggestion[] = [];
		
		for (const path of methodPaths) {
			const score = this.calculateScore(path, normalizedInput);
			if (score > 0) {
				suggestions.push({ path, score });
			}
		}
		
		// Sort by score (descending), then by length (ascending)
		return suggestions
			.sort((a, b) => {
				if (a.score !== b.score) {
					return b.score - a.score;
				}
				return a.path.length - b.path.length;
			})
			.slice(0, maxResults);
	}
	
	/**
	 * Calculate relevance score for a path given the input
	 */
	private static calculateScore(path: string, normalizedInput: string): number {
		const normalizedPath = path.toLowerCase();
		
		// Exact match gets highest score
		if (normalizedPath === normalizedInput) {
			return 100;
		}
		
		// Handle cases where user input might be missing leading slash
		const inputWithSlash = normalizedInput.startsWith('/') ? normalizedInput : '/' + normalizedInput;
		const inputWithoutSlash = normalizedInput.startsWith('/') ? normalizedInput.substring(1) : normalizedInput;
		
		// Starts with input gets high score (check both with and without leading slash)
		if (normalizedPath.startsWith(normalizedInput) || 
			normalizedPath.startsWith(inputWithSlash) || 
			normalizedPath.startsWith(inputWithoutSlash)) {
			return 90;
		}
		
		// Contains input gets medium score (check all variants)
		if (normalizedPath.includes(normalizedInput) || 
			normalizedPath.includes(inputWithSlash) || 
			normalizedPath.includes(inputWithoutSlash)) {
			return 60;
		}
		
		// Check if input matches path segments (parts between /)
		const inputParts = normalizedInput.split('/').filter(part => part.length > 0);
		const pathParts = normalizedPath.split('/').filter(part => part.length > 0);
		
		let matchingParts = 0;
		let consecutiveMatches = 0;
		let maxConsecutiveMatches = 0;
		
		for (let i = 0; i < inputParts.length; i++) {
			const inputPart = inputParts[i];
			let found = false;
			
			for (let j = 0; j < pathParts.length; j++) {
				const pathPart = pathParts[j];
				
				if (pathPart.includes(inputPart) || inputPart.includes(pathPart)) {
					matchingParts++;
					consecutiveMatches++;
					found = true;
					break;
				}
			}
			
			if (!found) {
				maxConsecutiveMatches = Math.max(maxConsecutiveMatches, consecutiveMatches);
				consecutiveMatches = 0;
			}
		}
		
		maxConsecutiveMatches = Math.max(maxConsecutiveMatches, consecutiveMatches);
		
		if (matchingParts > 0) {
			// Score based on matching parts and consecutive matches
			return Math.min(50, (matchingParts * 10) + (maxConsecutiveMatches * 5));
		}
		
		// Fuzzy matching - check if characters appear in order
		let inputIndex = 0;
		for (let i = 0; i < normalizedPath.length && inputIndex < normalizedInput.length; i++) {
			if (normalizedPath[i] === normalizedInput[inputIndex]) {
				inputIndex++;
			}
		}
		
		if (inputIndex === normalizedInput.length) {
			return Math.max(10, 30 - (normalizedPath.length - normalizedInput.length));
		}
		
		return 0;
	}
	
	/**
	 * Get all available methods
	 */
	static getAvailableMethods(): string[] {
		return Object.keys(elasticsearchMethodsPaths);
	}
}