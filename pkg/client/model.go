package client

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	ApiDomain    string `json:"api_domain"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}
type EmployeeResponse struct {
	Response struct {
		Result  []map[string][]Employee `json:"result"`
		Message string                  `json:"message"`
		URI     string                  `json:"uri"`
		Status  int                     `json:"status"`
	} `json:"response"`
}

type SingleEmployeeResponse struct {
	Response struct {
		Result  []Employee `json:"result"`
		Message string     `json:"message"`
		URI     string     `json:"uri"`
		Status  int        `json:"status"`
	} `json:"response"`
}
type Employee struct {
	MiddleName      string `json:"Middle_Name"`
	EmailID         string `json:"EmailID"`
	CreatedTime     string `json:"CreatedTime"`
	EmployeeTypeId  string `json:"Employee_type.id"`
	DateOfBirth     string `json:"Date_of_birth"`
	AddedTime       string `json:"AddedTime"`
	Photo           string `json:"Photo"`
	MaritalStatus   string `json:"Marital_status"`
	Gender          string `json:"Gender"`
	ModifiedBy      string `json:"ModifiedBy"`
	ApprovalStatus  string `json:"ApprovalStatus"`
	Department      string `json:"Department"`
	LocationNameID  string `json:"LocationName.ID"`
	TabularSections struct {
		EducationDetails []struct {
			Specialization   string `json:"Specialization"`
			Degree           string `json:"Degree"`
			College          string `json:"College"`
			YearOfGraduation string `json:"Yearofgraduation"`
			TabularROWID     string `json:"tabular.ROWID"`
		} `json:"Education Details"`
		WorkExperience []struct {
			JobTitle        string `json:"Jobtitle"`
			Employer        string `json:"Employer"`
			RELEVANCE       string `json:"RELEVANCE"`
			PreviousJobDesc string `json:"Previous_JobDesc"`
			FromDate        string `json:"FromDate"`
			ToDate          string `json:"Todate"`
			TabularROWID    string `json:"tabular.ROWID"`
			RELEVANCEId     string `json:"RELEVANCE.id"`
		} `json:"Work Experience"`
		DependentDetails []struct {
		} `json:"Dependent Details"`
	} `json:"tabularSections"`
	AddedBy                     string `json:"AddedBy"`
	MobileCountryCode           string `json:"Mobile.country_code"`
	Tags                        string `json:"Tags"`
	ReportingTo                 string `json:"Reporting_To"`
	PhotoDownloadUrl            string `json:"Photo_downloadUrl"`
	SourceOfHireId              string `json:"Source_of_hire.id"`
	TotalExperienceDisplayValue string `json:"total_experience.displayValue"`
	Citizenship                 string `json:"Citizenship"`
	EmployeeStatus              string `json:"Employeestatus"`
	Role                        string `json:"Role"`
	Experience                  string `json:"Experience"`
	EmployeeType                string `json:"Employee_type"`
	AddedByID                   string `json:"AddedBy.ID"`
	SocialSecurityNumber        string `json:"Social_Security_Number"`
	RoleID                      string `json:"Role.ID"`
	LastName                    string `json:"LastName"`
	EmployeeID                  string `json:"EmployeeID"`
	ZUID                        string `json:"ZUID"`
	VeteranStatus               string `json:"Veteran_Status"`
	DateOfExit                  string `json:"Dateofexit"`
	PermanentAddress            string `json:"Permanent_Address"`
	OtherEmail                  string `json:"Other_Email"`
	LocationName                string `json:"LocationName"`
	WorkLocation                string `json:"Work_location"`
	PresentAddress              string `json:"Present_Address"`
	NickName                    string `json:"Nick_Name"`
	TotalExperience             string `json:"total_experience"`
	ModifiedTime                string `json:"ModifiedTime"`
	ReportingToMailID           string `json:"Reporting_To.MailID"`
	ZohoID                      int64  `json:"Zoho_ID"`
	DesignationID               string `json:"Designation.ID"`
	SourceOfHire                string `json:"Source_of_hire"`
	Designation                 string `json:"Designation"`
	EthnicityId                 string `json:"Ethnicity.id"`
	MaritalStatusId             string `json:"Marital_status.id"`
	FirstName                   string `json:"FirstName"`
	AboutMe                     string `json:"AboutMe"`
	DateOfJoining               string `json:"Dateofjoining"`
	ExperienceDisplayValue      string `json:"Experience.displayValue"`
	Mobile                      string `json:"Mobile"`
	Extension                   string `json:"Extension"`
	ModifiedByID                string `json:"ModifiedBy.ID"`
	Ethnicity                   string `json:"Ethnicity"`
	ReportingToID               string `json:"Reporting_To.ID"`
	WorkPhone                   string `json:"Work_phone"`
	EmployeeStatusType          int    `json:"Employeestatus.type"`
	DepartmentID                string `json:"Department.ID"`
	PresentAddressChildValues   struct {
		CITY        string `json:"CITY"`
		COUNTRY     string `json:"COUNTRY"`
		STATE       string `json:"STATE"`
		ADDRESS1    string `json:"ADDRESS1"`
		PINCODE     string `json:"PINCODE"`
		ADDRESS2    string `json:"ADDRESS2"`
		STATECODE   string `json:"STATE_CODE"`
		COUNTRYCODE string `json:"COUNTRY_CODE"`
	} `json:"Present_Address.childValues"`
	Expertise string `json:"Expertise"`
}

type DepartmentResponse struct {
	Response struct {
		Result  []map[string][]Department `json:"result"`
		Message string                    `json:"message"`
		URI     string                    `json:"uri"`
		Status  int                       `json:"status"`
	} `json:"response"`
}

type SingleDepartmentResponse struct {
	Response struct {
		Result  []Department `json:"result"`
		Message string       `json:"message"`
		URI     string       `json:"uri"`
		Status  int          `json:"status"`
	} `json:"response"`
}

type Department struct {
	CreatedTime        string `json:"CreatedTime"`
	DepartmentLeadMail string `json:"Department_Lead.MailID"`
	AddedTime          string `json:"AddedTime"`
	DepartmentLead     string `json:"Department_Lead"`
	ModifiedBy         string `json:"ModifiedBy"`
	ApprovalStatus     string `json:"ApprovalStatus"`
	ModifiedByID       string `json:"ModifiedBy.ID"`
	Department         string `json:"Department"`
	DepartmentLeadID   string `json:"Department_Lead.ID"`
	ParentDepartmentID string `json:"Parent_Department.ID"`
	ModifiedTime       string `json:"ModifiedTime"`
	ZohoID             int    `json:"Zoho_ID"`
	AddedByID          string `json:"AddedBy.ID"`
	ParentDepartment   string `json:"Parent_Department"`
	AddedBy            string `json:"AddedBy"`
	MailAlias          string `json:"Mail_Alias"`
}
