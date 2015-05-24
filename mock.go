package main

var users map[int]contactList
var mockMessageId = 1000
func init () {


	users = make(map[int]contactList)

	users[101] = contactList{
		contactData{"ivan", 102, true},
		contactData{"petr", 103, true},
		contactData{"vladimir", 104, true},
	}

	users[102] = contactList{
		contactData{"huilo", 101, true},
		contactData{"petr", 103, true},
		contactData{"vladimir", 104, true},
	}

	users[103] = contactList{
		contactData{"ivan", 102, true},
		contactData{"huilo", 101, true},
		contactData{"vladimir", 104, true},
	}

	users[104] = contactList{
		contactData{"ivan", 102, true},
		contactData{"petr", 103, true},
		contactData{"huilo", 101, true},
	}
}
