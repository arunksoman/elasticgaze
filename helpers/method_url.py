import json

def extract_methods_urls_from_openapi(openapi_file_path, output_file_path):
    """
    Extract HTTP methods and URLs from OpenAPI 3.0 spec and create a JSON file
    in the format: {"method": [list of urls]}
    """
    
    # Read the OpenAPI spec file
    with open(openapi_file_path, 'r', encoding='utf-8') as f:
        openapi_spec = json.load(f)
    
    # Dictionary to store method -> URLs mapping
    method_url_mapping = {}
    
    # Extract paths from the OpenAPI spec
    paths = openapi_spec.get('paths', {})
    
    # Iterate through each path and its methods
    for path, path_data in paths.items():
        for method in path_data.keys():
            # Only consider HTTP methods (skip other keys like 'parameters', 'summary', etc.)
            if method.lower() in ['get', 'post', 'put', 'delete', 'head', 'options', 'patch', 'trace']:
                method_upper = method.upper()
                
                # Initialize the method key if it doesn't exist
                if method_upper not in method_url_mapping:
                    method_url_mapping[method_upper] = []
                
                # Add the path to the method's URL list
                method_url_mapping[method_upper].append(path)
    
    # Sort URLs for each method for better readability
    for method in method_url_mapping:
        method_url_mapping[method].sort()
    
    # Write the result to the output file
    with open(output_file_path, 'w', encoding='utf-8') as f:
        json.dump(method_url_mapping, f, indent=2, ensure_ascii=False)
    
    # Print summary
    total_endpoints = sum(len(urls) for urls in method_url_mapping.values())
    print(f"Successfully extracted {total_endpoints} endpoints across {len(method_url_mapping)} HTTP methods")
    print(f"Methods found: {', '.join(sorted(method_url_mapping.keys()))}")
    print(f"Output saved to: {output_file_path}")
    
    return method_url_mapping

if __name__ == "__main__":
    # File paths
    openapi_file = r"c:\Users\LENOVO\Downloads\elasticsearch-openapi-source.json"
    output_file = r"d:\desktop_apps\elasticgaze\helpers\elasticsearch_methods_paths.json"
    
    # Extract methods and URLs
    result = extract_methods_urls_from_openapi(openapi_file, output_file)
    
    # Display some statistics
    print("\nStatistics by method:")
    for method, urls in sorted(result.items()):
        print(f"  {method}: {len(urls)} endpoints")
