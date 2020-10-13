import React from 'react';
import './App.css';

var host = "http://192.168.128.146:9427/";

class GeneralInfoPage extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      process : "",
      sites : [],
    };
  }
  componentDidMount() {
    let url = host + "info";
    let xhr = new XMLHttpRequest();
    xhr.addEventListener('load', () => {
      this.setState({
        process : JSON.parse(xhr.responseText.replace(/\t/g, ' ')).currentProcess,
        sites : JSON.parse(xhr.responseText.replace(/\t/g, ' ')).siteNames,
      });
    })
    xhr.open('GET', url);
    xhr.send();
  }
  redirect(destination) {
    destination = '/' + destination + '/';
    return () => this.props.history.push(destination);
  }
  render() {
    let items = [];
      items.push(<p>Process : </p>)
    if (this.state.process === "") {
      items.push(<button className="empty-process" onClick={this.redirect("process")}>empty</button>)
      items.push(<br/>)
    } else {
      items.push(<button className="running-process" onClick={this.redirect('process')}>{this.state.process}</button>)
      items.push(<br/>)
    }
    for (let i in this.state.sites) {
      items.push(<button className="site" onClick={this.redirect(this.state.sites[i])}>{this.state.sites[i]}</button>);
    }
    return (
    <div className="container">
      <header><p>Book</p></header>
      <div className="Page">
        {items}
      </div>
    </div>
    );
  }
}

class ProcessInfoPage extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      process : "unknown",
    }
  }
  componentDidMount() {
    let url = host + "info";
    let xhr = new XMLHttpRequest();
    xhr.addEventListener('load', () => {
      this.setState({
        process : JSON.parse(xhr.responseText.replace(/\t/g, ' ')).currentProcess
      });
    })
    xhr.open('GET', url);
    xhr.send();
  }
  redirect(destination) {
    return () => this.props.history.push(destination);
  }
  start(processName) {
    return () => {
      let url = host + "start?operation="+processName.replace(' ', '');
      let xhr = new XMLHttpRequest();
      xhr.open('GET', url);
      xhr.send();
    }
  }
  render() {
    let process = (this.state.process === "") ? (<button className="empty-process">empty</button>) : (<button className="running-process">{this.state.process}</button>)
    let items = []
    let processes = ["Update", "Explore", "Download", "Error", "Check", "Check End", "Backup", "Fix"]
    for (let i in processes) {
      if (this.state.process === "") {
        items.push(<button className="able-process" onClick={this.start(processes[i])}>{processes[i]}</button>)
      } else {
        items.push(<button className="able-process" disabled>{processes[i]}</button>)
      }
      items.push(<br/>)
    }
    return (
      <div className="container">
        <header><p onClick={this.redirect("/")}>Book</p> > <p>Process</p></header>
        <div className="page">
          <p>Process : </p>{process}<br/>
          <div className="scroller">
            {items}
          </div>
          <br/>
          <button onClick={this.redirect("/logs/")}>Logs</button>
        </div>
      </div>
    )
  }
}

class LogsInfoPage extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      datetime : "",
      logs : []
    }
  }
  componentDidMount() {
    let url = host + "process";
    let xhr = new XMLHttpRequest();
    xhr.addEventListener('load', () => {
      try {
        let obj = JSON.parse(xhr.responseText.replace(/\t/g, ' ').replace(/="/g, '=\'').replace(/" (\w)/g, '\' $1').replace(/">/g, '\'>'));
        this.setState({
          datetime : obj.time,
          logs : obj.logs,
        });
      } catch (e){
        this.setState({
	  datetime : "unknown",
	  logs : xhr.responseText.replace(/"/g, "").split(','),
	});
      }
    });
    xhr.open('GET', url);
    xhr.send();
  }
  redirect(destination) {
    return () => this.props.history.push(destination);
  }
  render() {
    let logs = []
    for (let i in this.state.logs) {
      logs.push(<p>{this.state.logs[i]}</p>)
      logs.push(<br/>)
    }
    console.log(logs)
    return (
      <div className="container">
        <header><p onClick={this.redirect("/")}>Book</p> > <p onClick={this.redirect("/process/")}>Process</p> > <p>Logs</p></header>
        <div className="page">
          <p>Datetime : {this.state.datetime}</p>
          <br/><br/>
          <hr/>
          <div className="scroller" style={{height: "70vh"}}>
            {logs}
          </div>
        </div>
      </div>
    )
  }
}

class SiteInfoPage extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      name : props.match.params.name,
      formValues : {
        title : "",
        writer : ""
      }
    };
    this.search = this.search.bind(this)
  }
  componentDidMount() {
    let url = host + "info/"+this.state.name;
    let xhr = new XMLHttpRequest();
    xhr.addEventListener('load', () => {
      let obj = JSON.parse(xhr.responseText.replace(/\t/g, ' '))
      this.setState({
        bookCount : obj.bookCount,
        errorCount : obj.errorCount,
        endCount : obj.endCount,
        downloadCount : obj.downloadCount,
        bookRecordCount : obj.bookRecordCount,
        errorRecordCount : obj.errorRecordCount,
        endRecordCount : obj.endRecordCount,
        downloadRecordCount : obj.downloadRecordCount,
        readCount : obj.readCount,
        maxNum : obj.maxid
      });
    })
    xhr.open('GET', url);
    xhr.send();
  }
  redirect(destination) {
    return () => this.props.history.push(destination);
  }
  search() {
    let url = "/" + this.state.name + "/search?title=" + this.state.formValue.title + "&writer=" + this.state.formValue.writer
    console.log(url)
    this.redirect(url)();
    return false
  }
  handleChange(event) {
    event.preventDefault();
    let formValues = this.state.formValues;
    let name = event.target.name;
    let value = event.target.value;

    formValues[name] = value;

    this.setState({
      formValue : formValues
    });
}
  render() {
    return (
      <div className="container">
        <header><p onClick={this.redirect("/")}>Book</p> > <p>Site</p> > <p>{this.state.name}</p></header>
        <div className="page">
          <p>Book Count : {this.state.bookCount + this.state.errorCount} ({this.state.bookCount} + {this.state.errorCount})</p><br/>
          <p>Record Count : {this.state.bookRecordCount + this.state.errorRecordCount} ({this.state.bookRecordCount} + {this.state.errorRecordCount})</p><br/>
          <p>End Count : {this.state.endCount} ({this.state.endRecordCount})</p><br/>
          <p>Download Count : {this.state.downloadCount} ({this.state.downloadRecordCount})</p><br/>
          <p>Read Count : {this.state.readCount}</p><br/>
          <p>Max ID : {this.state.maxNum}</p><br/>
          <form onSubmit={this.search} autoComplete="off">
            <p>Search</p><br/>
            <hr/>
            <label for="title">Title : </label>
            <input type="text" id="title" name="title" placeholder="Title" value={this.state.searchTitle} onChange={this.handleChange.bind(this)}/><br/>
            <label for="writer">Writer : </label>
            <input type="text" id="writer" name="writer" placeholder="Writer name" value={this.state.searchWriter} onChange={this.handleChange.bind(this)}/><br/>
            <input type="submit"></input>
          </form>
        </div>
      </div>
    )
  }
}

class SearchBookPage extends React.Component {
  constructor(props) {
    super(props)
    let uri = new URLSearchParams(this.props.location.search)
    this.state = {
      name : props.match.params.name,
      page : 0,
      title : uri.get('title'),
      writer : uri.get('writer'),
      books : []
    }
    this.search = this.search.bind(this);
    this.nextPage = this.nextPage.bind(this);
    this.lastPage = this.lastPage.bind(this);
  }
  componentDidMount(page=0) {
    let url = 'http://192.168.128.146:9427/search/'+this.state.name+"?title="+
      this.state.title+"&writer="+this.state.writer+"&page="+page;
    let xhr = new XMLHttpRequest();
    xhr.addEventListener('load', () => {
      this.setState({
        books : JSON.parse(xhr.responseText.replace(/\t/g, ' ')).books
      });
    })
    xhr.open('GET', url);
    xhr.send();
  }
  redirect(destination) {
    return () => this.props.history.push(destination);
  }
  search() {
    let url = "/" + this.state.name + "/search?title=" + this.state.formValue.title + "&writer=" + this.state.formValue.writer
    console.log(url)
    this.redirect(url)();
    return false
  }
  nextPage() {
    if (this.state.books.length < 20) {
      return
    }
    this.setState({page : this.state.page + 1});
    this.componentDidMount(this.state.page + 1);
  }
  lastPage() {
    if (this.state.page === 0) {
      return
    }
    this.setState({page : this.state.page - 1});
    this.componentDidMount(this.state.page - 1);
  }
  render() {
    let books = []
    for (let i in this.state.books) {
      let book = this.state.books[i]
      books.push(<table className="book" onClick={this.redirect("/"+this.state.name+"/book/"+book.num)}>
        <tbody>
          <tr><td className="title-info">{book.title}</td><td className="writer-info">{book.writer}</td></tr>
          <tr><td className="date-info">{book.update}</td><td className="chapter-info">{book.chapter}</td></tr>
        </tbody>
      </table>)
      books.push(<hr/>)
    }
    if (books.length === 0) {
      books.push(<p>No matched result found</p>);
    }
    let nextAvail = false
    let lastAvail = false
    if (this.state.books.length < 20) {
      nextAvail = true;
    }
    if (this.state.page === 0) {
      lastAvail = true;
    }
    return (
      <div className="container">
        <header><p onClick={this.redirect("/")}>Book</p> > <p>Search</p></header>
        <div className="page">
          <form onSubmit={this.search} autoComplete="off">
            <label for="title">Title : </label>
            <input type="text" id="title" name="title" placeholder="Title"value={this.state.searchTitle}/><br/>
            <label for="writer">Writer : </label>
            <input type="text" id="writer" name="writer" placeholder="Writer name" value={this.state.searchWriter}/><br/>
            <input type="submit"></input>
          </form>
          <div className="page">
            <button className="lastPage" onClick={this.lastPage} disabled={lastAvail}>&lt;</button>
            <p className="page">{this.state.page}</p>
            <button className="nextPage" onClick={this.nextPage} disabled={nextAvail}>&gt;</button>
          </div>
          <div className="scroller">
            {books}
          </div>
        </div>
      </div>
    )
  }
}

class BookInfoPage extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      name : props.match.params.name,
      num : props.match.params.num,
      title : "unknown",
      writer : "unknown",
      update : "unknown",
      chapter : "unknown"
    }
  }
  componentDidMount() {
    let url = 'http://192.168.128.146:9427/info/' + this.state.name + '/' + this.state.num;
    let xhr = new XMLHttpRequest();
    xhr.addEventListener('load', () => {
      let book = JSON.parse(xhr.responseText.replace(/\t/g, ' '));
      this.setState({
        title : book.title,
        writer : book.writer,
        type : book.type,
        update : book.update,
        chapter : book.chapter,
        version : book.version,
        download : book.download

      });
    })
    xhr.open('GET', url);
    xhr.send();

  }
  redirect(destination) {
    return () => this.props.history.push(destination);
  }
  download() {
    return () => window.open('http://192.168.128.146:9427/download/'+this.state.name+"/"+this.state.num);
  }
  search() {
    return () => window.open('https://www.google.com/search?q='+this.state.name+"+"+this.state.title)
  }
  render() {
    let button = null;
    if (this.state.download === true) {
      button = <button className="download" onClick={this.download()}>Download</button>
    } else {
      button = <button className="search" onClick={this.search()}>Search Online</button>
    }
    return (
      <div className="container">
        <header><p onClick={this.redirect("/")}>Book</p> > <p onClick={this.redirect("/"+this.state.name+"/")}>{this.state.name}</p> > <p>{this.state.title}</p></header>
        <div className="page">
          <p>Title : {this.state.title}</p><br/>
          <p>Writer : {this.state.writer}</p><br/>
          <p>Type : {this.state.type}</p><br/>
          <p>Update : {this.state.update}</p><br/>
          <p>Chapter : {this.state.chapter}</p><br/>
          <p>Version : {this.state.version}</p><br/>
          {button}
        </div>
      </div>
    )
  }
}

export {
  GeneralInfoPage,
  ProcessInfoPage,
  LogsInfoPage,
  SiteInfoPage,
  SearchBookPage,
  BookInfoPage
};
