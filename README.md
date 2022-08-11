# TN-Plates

Quickly determine if a custom license plate is available in Tennessee.

### usage

```
# plate contains a word that is forbidden in this case 'hell'
% ./plates -p "hello"
hello: Plate number hello contains objectionable words

# plate is already in use
% ./plates -p "irun" 
irun: Rights are active for this plate for the next 597 days

# plate is available for a custom license plate
% ./plates -p "irunz"
irunz: available

# supplying multiple plate requests is also supported
% ./plates -p "hello" -p "irun" -p "irunz"
hello: Plate number hello contains objectionable words
irun:  Rights are active for this plate for the next 597 days
irunz: available

```
