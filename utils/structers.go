package utils

type Mchanges struct {
	Changed []string `json:"changed"`
	MTime   []string `json:"m_time"`
}

type Brecord struct {
	Category string   `json:"category"`
	FDirs     []string `json:"f_dir"`
	LMod     string   `json:"l_mod"`
	Changes  Mchanges `json:"changes"`
}


