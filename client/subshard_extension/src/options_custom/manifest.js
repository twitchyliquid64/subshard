
this.manifest = {
    "name": "Subshard guard",
    "icon": "../../icons/chromeball_google_chrome_poke_by_azerik92-d4c31vz.png",
    "settings": [
        {
            "tab": "Information",
            "group": "Overview",
            "name": "myDescription",
            "type": "description",
            "text": "Welcome to Subshard! If you can load a page at <a href=\"http://subshard\">http://subshard/</a>, you have set everything up correctly."
        },
        {
            "tab": "Information",
            "group": "Overview",
            "name": "myDescription2",
            "type": "description",
            "text": "This extension exists to help protect you online. It simply modifies some headers to make it harder to track what sites you have been to, and prevents resource leakage across sensitive domains."
        },
        {
            "tab": "Settings",
            "group": "General",
            "name": "cross_resources_disabled",
            "type": "checkbox",
            "label": "Prevent sites from requesting resources from different domains (disable CORS)"
        },
        {
            "tab": "Settings",
            "group": "General",
            "name": "sensitive_protection_disabled",
            "type": "checkbox",
            "label": "Allow sensitive sites to request resources from a different domain (CORS)"
        },
    ],
    "alignment": [
    ]
};
