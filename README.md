# ctx
Switching context in your shell is usually a nightmare. Say no more to copy and pasting values from your Keepass database. CTX is used to read Keepass entries and store them in your environment variables. It does so by reading a root group and searching for entries based on their title/environment (the title/environment is supplied with the `-e` option). It will print out every line in the `notes` section to STDOUT. This output can be used in combination with `eval $()` in order to set your environment.

## Prereqs
Specify the location of your Keepass database by setting the environment variable `CTX_VAR_db_location`. To open the Keepass database, set the environment variable `CTX_VAR_secret`.

## Usage
By default the Keepass root group called `Azure` will be used. Can be overwritten by specifying it with the option `-g groupname`. When executing the program, it will only output the environment values to STDOUT. It won't be able to set the environment variables in your shell since you can't change the environment of the parent process.

However, we can do a neat little trick by using `eval $()` to export the output of the program in your shell. Example:
```bash
eval $(ctx -e maz000-p)
```

### Example configuration keepass
```
keepassdatabasename
|
-- Azure
 |
 -- Title: maz000-p
    Notes:
        TF_VAR_tenant_id=<TENANTID>
        TF_VAR_subscription_id=<SUBSCRIPTIONID>
```

Output:
```bash
$ ctx -e maz000-p -g Azure
 export TF_VAR_tenant_id=<TENANTID>
 export TF_VAR_subscription_id=<SUBSCRIPTIONID>
```


# TODO
- Check state of current env var's
