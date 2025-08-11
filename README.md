1. auth context 
2. repo -> file read & write 
3. serice & repo ( interface )
4. jwt 
5. enum


REPO 

InitServices 

-> user, auth, product, category, lender
-> servcies 

userrepo (user.json, lender.json)
authService(userrepo)


userrepo (user.json, lenderRepo)


userrepo {
    userFIle
    lenderRepo,
    categoryRepo
}


<!--  Read & WRITE -->

Init -> .NewUserFileRepo()

type user_repo struct {
    userFile
    users[] <- data
    lenderRepo
}

become a lender -> users[]

users[] -> modify 

userrepo.save(array) write() 

updated data - write 


<!-- REMAINING THINGS -->

1. Feedback -> Done
2. Remaining Checks all over the project 
3. Category -> In Progress
4. Society ->  Done
5.  


1. Features Color
3. Create Account -> society input
4. Multiple order status input ( comma separated )
5. Get all pending return request ( add in features list )
