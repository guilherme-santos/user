# User Service

# Assumptions & Limitations

1. This API aims to be used by a admins, not for the user itself.
   * My concern here is regarding the API structure and password. I would have a endpoint `/v1/user` and provide a mechanism to change the password, requesting the current and the new password. For admins purpose in general you shouldn't have the current password, you'd be allowed to just replace it.
2. I'm using soft-delete for user removal. As far as GDPR is concerned, the soft-delete could also anonymise the personal data.
3. I've created a dummy cache layer just to show how I could extend the current codebase adding more functionality if necessary.
4. I've created a dummy publisher just to make it clear where I would have implemented.
5. The cursor implementation can be used only to move forward in the list sorted by id, it'll be useful for infinite scroll.
6. I'm using `github.com/rs/xid` to generate the ids which uses the Mongo Object ID algorithm. Main reason here is that I can sort by id and the result will be order of creation as it has a timestamp component.
7. I'm using `https://github.com/golang/mock` mock generation.
