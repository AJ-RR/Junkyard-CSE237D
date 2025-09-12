
echo "🔧 GitHub Actions Runner Setup Wizard"
printf "Enter your GitHub org URL (e.g. https://github.com/my-org/my-repo): "
read GITHUB_CONFIG_URL

printf "Create a Github Personal Access Token (classic)\n"
printf "Note: Ensure that your token as scope admin:org \n"
printf "Enter your GitHub Personal Access Token: "
read GITHUB_PAT

helm install arc \
    --namespace "greengrader-systems" \
    --create-namespace \
    oci://ghcr.io/actions/actions-runner-controller-charts/gha-runner-scale-set-controller

helm install "greengrader" \
    --namespace "greengrader-runners" \
    --create-namespace \
    --set githubConfigUrl="$GITHUB_CONFIG_URL" \
    --set githubConfigSecret.github_token="$GITHUB_PAT" \
    oci://ghcr.io/actions/actions-runner-controller-charts/gha-runner-scale-set
