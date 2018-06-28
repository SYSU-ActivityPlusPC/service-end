package dao

// Activity
func (act *ActivityInfo) Insert() {}

func (act *ActivityInfo) Delete() {}

func (act *ActivityInfo) Update() {}

func (act *ActivityInfo) Get() {}

func (act *ActivityInfo) GetListByClubID() {}

func (act *ActivityInfo) GetListByVerifiedStatus() {}

// Apply
func (apply *ActApplyInfo) GetApplyByActID() {}

func (apply *ActApplyInfo) Delete() {}

// PC user
func (user *PCUser) Get() {}

func (user *PCUser) Update() {}

func (user *PCUser) ListUserByType() {}

func (user *PCUser) Insert() {}

// Message
func (msg *Message) Insert() {}

func (msg *Message) List() {}

func (msg *Message) Get() {}

func (msg *Message) Delete() {}
