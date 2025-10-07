line = "example_line"
toadd = """
        //js
        <li>
            <a href="#{line}">{line}</a>
        </li>
        """
        
toadd = toadd.format(line=line)
print(toadd)