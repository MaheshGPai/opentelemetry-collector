# opentelemetry-collector
This repo is a mirror of https://github.com/open-telemetry/opentelemetry-collector. All Intuit specific customizations 
are maintained in this repo. The changes are kept in `intuit` branch

## Step To Release New Version

1. Clone the opensource repo
```shell
git clone git@github.com:opentelemetry/opentelemetry-collector.git
```

2. Go into the repo and link the intuit repo
```shell
git remote add intuit_remote git@github.intuit.com:opentracing-collector/opentelemetry-collector.git
```

3. Checkout the intuit branch from the custom/target repo
```shell
git switch -c intuit intuit_remote/intuit
```

4. Branch out for a new release (e.g) v0.103.0
```shell
git checkout intuit -b intuit_v0.103.0
```

4. Rebase with the release tag that we are upgrading to (e.g) v0.103.0
```shell
git rebase v0.103.0
```

5. Commit and push the changes
```shell
git commit -m "Intuit specific fixes"
git push intuit intuit_v0.103.0
```

6. Create and push new tags
```shell
sh createAndPushTags.sh v0.103.0-1
```