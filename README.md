[![GoDoc](https://godoc.org/gitlab.com/jannickfahlbusch/agenda_lgo?status.svg)](https://godoc.org/gitlab.com/jannickfahlbusch/agenda_lgo)

This is a simple program to download all of your salary statements from "Agenda: Lohn- und Gehaltsdokumente".

---

To run it, you need to specify the parameters "-a" and "-o": 

* -a: Path to the authentication-file  
        This is a JSON file with the following structure:
        
	```
	{
        	"Email": "",
            "Password": ""
	}
	```

* -o: Path to the directory, where all PDFs should be stored




