package trackers

// Torrent defines a local torrent object, a smaller subset of original data.
type Torrent struct {
	ID       int        // torrent ID in torrent-client
	Name     string     // torrent name
	Trackers []*Tracker // torrent tracker's, based on local type
}
