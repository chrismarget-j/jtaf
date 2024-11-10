# To Do:
- add git commit and git tag variables
- collect hash of yang data
- calculate provider name
- ~~clone url.QueryEscape so that it works for our characters~~
- ~~semaphore for client transactions~~
- eliminate path attribute for top-level resources
- introduce function to convert XML config attributes to terraform attributes
  - Figure out whether any yang leafs begin with non-alpha characters
  - `native-inner-vlan-id` -> `native_inner_vlan_id`