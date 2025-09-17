import os
import requests

# Load configuration from environment variables
GITLAB_HOST = os.getenv('GITLAB_HOST')
GROUP_ID = os.getenv('GROUP_ID')
ACCESS_TOKEN = os.getenv('ACCESS_TOKEN')
SEARCH_VALUE = os.getenv('SEARCH_VALUE')

# Validate inputs
if not all([GITLAB_HOST, GROUP_ID, ACCESS_TOKEN, SEARCH_VALUE]):
    raise ValueError("Missing one or more required environment variables.")

# Set headers for GitLab API
headers = {
    'Private-Token': ACCESS_TOKEN
}

def get_group_projects():
    """Fetch all projects in the GitLab group."""
    projects = []
    page = 1
    while True:
        url = f"{GITLAB_HOST}/api/v4/groups/{GROUP_ID}/projects?per_page=100&page={page}"
        response = requests.get(url, headers=headers)
        if response.status_code != 200:
            print(f"Error fetching projects: {response.status_code}")
            break
        data = response.json()
        if not data:
            break
        projects.extend(data)
        page += 1
    return projects

def check_project_variables(project_id):
    """Check CI/CD variables for a specific value."""
    url = f"{GITLAB_HOST}/api/v4/projects/{project_id}/variables"
    response = requests.get(url, headers=headers)
    if response.status_code != 200:
        print(f"Error fetching variables for project {project_id}: {response.status_code}")
        return []
    variables = response.json()
    return [var for var in variables if SEARCH_VALUE in var.get('value', '')]

def main():
    projects = get_group_projects()
    for project in projects:
        project_id = project['id']
        project_name = project['name']
        matching_vars = check_project_variables(project_id)
        if matching_vars:
            print(f"\nProject: {project_name} (ID: {project_id})")
            for var in matching_vars:
                print(f"  Variable Key: {var['key']}, Value: {var['value']}")
        else:
            print(f"  No keys for Project: {project_name} (ID: {project_id})")


if __name__ == "__main__":
    main()
