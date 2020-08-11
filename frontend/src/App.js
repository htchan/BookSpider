import React from 'react';
import './App.css';

class GeneralInfoPage extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      process : "",
      sites : [],
    };
  }
  componentDidMount() {
    let url = "http://192.168.128.146:9427/info";
    let xhr = new XMLHttpRequest();
    xhr.addEventListener('load', () => {
      this.setState({
        process : JSON.parse(xhr.responseText).currentProcess,
        sites : JSON.parse(xhr.responseText).siteNames,
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
      items.push(<a>Process : </a>)
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
      <header><a>Book</a></header>
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
    let url = "http://192.168.128.146:9427/info";
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
      let url = "http://192.168.128.146:9427/start?operation="+processName;
      let xhr = new XMLHttpRequest();
      xhr.open('GET', url);
      xhr.send();
    }
  }
  render() {
    let process = (this.state.process === "") ? (<button className="empty-process">empty</button>) : (<button className="running-process">{this.state.process}</button>)
    let items = []
    let processes = ["Update", "Explore", "Download", "Error", "Check", "Backup", "Fix"]
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
        <header><a onClick={this.redirect("/")}>Book</a> > <a>Process</a></header>
        <div className="page">
          <a>Process : </a>{process}<br/>
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
    let url = "http://192.168.128.146:9427/process";
    let xhr = new XMLHttpRequest();
    xhr.addEventListener('load', () => {
      this.setState({
        datetime : JSON.parse(xhr.responseText.replace(/\t/g, ' ')).time,
        logs : JSON.parse(xhr.responseText.replace(/\t/g, ' ')).logs,
      });
    })
    xhr.open('GET', url);
    xhr.send();
  }
  redirect(destination) {
    return () => this.props.history.push(destination);
  }
  render() {
    let logs = []
    for (let i in this.state.logs) {
      logs.push(<a>{this.state.logs[i]}</a>)
      logs.push(<br/>)
    }
    console.log(logs)
    return (
      <div className="container">
        <header><a onClick={this.redirect("/")}>Book</a> > <a onClick={this.redirect("/process/")}>Process</a> > <a>Logs</a></header>
        <div className="page">
          <a>Datetime : {this.state.datetime}</a>
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
      name : props.match.params.name
    };
  }
  componentDidMount() {
    let url = "http://192.168.128.146:9427/info/"+this.state.name;
    let xhr = new XMLHttpRequest();
    xhr.addEventListener('load', () => {
      this.setState({
        bookCount : JSON.parse(xhr.responseText).bookCount,
        errorCount : JSON.parse(xhr.responseText).errorCount,
        totalCount : JSON.parse(xhr.responseText).totalCount,
        bookVersionCount : JSON.parse(xhr.responseText).bookVersionCount,
        endCount : JSON.parse(xhr.responseText).endCount,
        downloadCount : JSON.parse(xhr.responseText).downloadCount,
        readCount : JSON.parse(xhr.responseText).readCount,
        maxNum : JSON.parse(xhr.responseText).maxNum
      });
    })
    xhr.open('GET', url);
    xhr.send();
  }
  redirect(destination) {
    return () => this.props.history.push(destination);
  }
  render() {
    return (
      <div className="container">
        <header><a onClick={this.redirect("/")}>Book</a> > <a>Site</a> > <a>{this.state.name}</a></header>
        <div className="page">
          <a>Book Count : {this.state.bookCount}</a><br/>
          <a>Error Count : {this.state.errorCount}</a><br/>
          <a>Total Count : {this.state.totalCount}</a><br/>
          <a>Book Version Count : {this.state.bookVersionCount}</a><br/>
          <a>End Count : {this.state.endCount}</a><br/>
          <a>Download Count : {this.state.downloadCount}</a><br/>
          <a>Read Count : {this.state.readCount}</a><br/>
          <a>Max ID : {this.state.maxNum}</a><br/>
          <form action={"/"+this.state.name+"/search"} autoComplete="off">
            <a>Search</a><br/>
            <hr/>
            <label for="title">Title : </label>
            <input type="text" id="title" name="title" placeholder="Title" value={this.state.searchTitle}/><br/>
            <label for="writer">Writer : </label>
            <input type="text" id="writer" name="writer" placeholder="Writer name" value={this.state.searchWriter}/><br/>
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
      title : uri.get('title'),
      writer : uri.get('writer')
    }
  }
  componentDidMount() {
    let url = 'http://192.168.128.146:9427/search/'+this.state.name+"?title="+
      this.state.title+"&writer="+this.state.writer;
    let xhr = new XMLHttpRequest();
    xhr.addEventListener('load', () => {
      this.setState({
        books : JSON.parse(xhr.responseText).books
      });
    })
    xhr.open('GET', url);
    xhr.send();
  }
  redirect(destination) {
    return () => this.props.history.push(destination);
  }
  render() {
    let books = []
    for (let i in this.state.books) {
      let book = this.state.books[i]
      books.push(<table className="book" onClick={this.redirect("/"+this.state.name+"/book/"+book.num)}>
        <tr><td className="title-info">{book.title}</td><td className="writer-info">{book.writer}</td></tr>
        <tr><td className="date-info">{book.update}</td><td className="chapter-info">{book.chapter}</td></tr>
      </table>)
      books.push(<hr/>)
    }
    if (books.length === 0) {
      books.push(<a>No matched result found</a>)
    }
    return (
      <div className="container">
        <header><a onClick={this.redirect("/")}>Book</a> > <a>Search</a></header>
        <div className="page">
          <form action={"/"+this.state.name+"/search"} autoComplete="off">
            <label for="title">Title : </label>
            <input type="text" id="title" name="title" placeholder="Title"value={this.state.searchTitle}/><br/>
            <label for="writer">Writer : </label>
            <input type="text" id="writer" name="writer" placeholder="Writer name" value={this.state.searchWriter}/><br/>
            <input type="submit"></input>
          </form>
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
      let book = JSON.parse(xhr.responseText);
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
        <header><a onClick={this.redirect("/")}>Book</a> > <a onClick={this.redirect("/"+this.state.name+"/")}>{this.state.name}</a> > <a>{this.state.title}</a></header>
        <div className="page">
          <a>Title : {this.state.title}</a><br/>
          <a>Writer : {this.state.writer}</a><br/>
          <a>Type : {this.state.type}</a><br/>
          <a>Update : {this.state.update}</a><br/>
          <a>Chapter : {this.state.chapter}</a><br/>
          <a>Version : {this.state.version}</a><br/>
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
